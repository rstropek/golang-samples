package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	v "github.com/auth0/go-jwt-middleware/v2/validator"
)

type CustomClaimsExample struct {
	Name       string `json:"name"`
	FamilyName string `json:"family_name"`
}

// Validate does nothing for this example.
func (c *CustomClaimsExample) Validate(ctx context.Context) error {
	return nil
}

func NewJwtMiddleware(azureTenantId string, requiredScopes []string) *jwtmiddleware.JWTMiddleware {
	issuerURL, err := url.Parse(fmt.Sprintf("https://login.microsoftonline.com/%s/", azureTenantId))
	if err != nil {
		log.Fatalf("failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	customClaims := &CustomClaimsExample{}
	jwtValidator, _ := v.New(
		provider.KeyFunc,
		"RS256",
		fmt.Sprintf("https://sts.windows.net/%s/", azureTenantId),
		requiredScopes,
		v.WithCustomClaims(func () v.CustomClaims { return customClaims }))
	return jwtmiddleware.New(jwtValidator.ValidateToken)
}

func ClaimsHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*v.ValidatedClaims)

	payload, err := json.Marshal(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
