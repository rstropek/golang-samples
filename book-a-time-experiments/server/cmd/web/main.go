package main

import (
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	jwtmiddleware "github.com/ciehanski/go-jwt-middleware"
	oidc "github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rstropek/book-a-time/cmd/web/handlers"
	"github.com/rstropek/book-a-time/pkg/auth0"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
	"golang.org/x/oauth2"
)

type Response struct {
	Message string `json:"message"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	gob.Register(map[string]interface{}{})

	viper.SetDefault("Port", ":5000")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/book-a-time/")
	viper.AddConfigPath("$HOME/.book-a-time")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			infoLog.Print("No config file found")
		} else {
			errorLog.Fatalf("Error reading config file: %+v", err)
		}
	}

	pcr := auth0.PemCertificateOidcJwksReader{
		BaseOidcURL: viper.GetString("auth.oidcBaseURL"),
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: auth0.ValidationKeyGetter(viper.GetString("auth.aud"), viper.GetString("auth.iss"), pcr),
		Debug:               viper.GetBool("auth.debug"),
	})

	r := mux.NewRouter()
	srv := http.Server{
		Addr:    viper.GetString("Port"),
		Handler: r,
	}

	hndlrs := handlers.Handlers{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Store:    store,
	}

	r.HandleFunc("/callback", CallbackHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/user", hndlrs.UserHandler)
	r.HandleFunc("/", hndlrs.Home)

	// This route is always accessible
	r.Handle("/api/public", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fooValue := session.Values["foo"]
		if fooValue == nil {
			session.Values["foo"] = fmt.Sprintf("bar%d", rand.Intn(10000))
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		message := fmt.Sprintf("Hello from a public endpoint! You don't need to be authenticated to see this. %v", session.Values["foo"])
		responseJSON(message, w, http.StatusOK)
	}))

	// This route is only accessible if the user has a valid Access Token
	// We are chaining the jwtmiddleware middleware into the negroni handler function which will check
	// for a valid token.
	r.Handle("/api/private", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			message := "Hello from a private endpoint! You need to be authenticated to see this."
			responseJSON(message, w, http.StatusOK)
		}))))

	idleConnectionsClosed := make(chan struct{})
	go func() {
		// wait for interrupt signal from terminal or kubernetes
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Fatalf("HTTP server Shutdown: %v", err)
		}

		close(idleConnectionsClosed)
	}()

	infoLog.Printf("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnectionsClosed
}

func responseJSON(message string, w http.ResponseWriter, statusCode int) {
	response := Response{message}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://rainerdemo.eu.auth0.com/")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     "8N7biYQcBdvawlKv8Q7Jiv7GKvE7DmYm",
		ClientSecret: "9_u6QTR0HyqeC2ayeA3cZUrzeePEqebuinKkQTXALO-6hVPDU01PIpvOs_1_IW8e",
		RedirectURL:  "http://localhost:5000/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := authenticator.Config.Exchange(context.TODO(), code)
	if err != nil {
		log.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: "8N7biYQcBdvawlKv8Q7Jiv7GKvE7DmYm",
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = rawIDToken
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)

	session, err := store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
