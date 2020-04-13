package auth0

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/dgrijalva/jwt-go"
)

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type jwks struct {
	Keys []jsonWebKeys `json:"keys"`
}

// PemCertificateReader reads PEM certificate for given key identifier (kid)
type PemCertificateReader interface {
	GetPemCertificate(kid interface{}) (string, error)
}

// PemCertificateOidcJwksReader reads PEM certificate from OIDC jwks.json
type PemCertificateOidcJwksReader struct {
	HTTPClient  *http.Client
	BaseOidcURL string
}

// GetPemCertificate reads the PEM cert from a given OpenID Connect provider
func (pcr PemCertificateOidcJwksReader) GetPemCertificate(kid interface{}) (cert string, err error) {
	cert = ""

	jwksURL, err := buildJwksWellKnownURL(pcr.BaseOidcURL)
	if err != nil {
		return
	}

	httpClient := pcr.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Get(jwksURL)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return cert, fmt.Errorf("HTTP status code 404 when accessing %s", pcr.BaseOidcURL)
	}

	defer resp.Body.Close()

	var jwks = jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return
	}

	for k := range jwks.Keys {
		if kid == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err = errors.New("Unable to find appropriate key")
		return
	}

	return cert, nil
}

func buildJwksWellKnownURL(baseOidcURL string) (string, error) {
	u, err := url.Parse(baseOidcURL)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, ".well-known", "jwks.json")
	return u.String(), nil
}

// ValidationKeyGetter returns a function that can be used with github.com/ciehanski/go-jwt-middleware
func ValidationKeyGetter(expectedAudience string, expectedIssuer string, pcr PemCertificateReader) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(expectedAudience, false)
		if !checkAud {
			return token, errors.New("Invalid audience")
		}

		// Verify 'iss' claim
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(expectedIssuer, false)
		if !checkIss {
			return token, errors.New("Invalid issuer")
		}

		cert, err := pcr.GetPemCertificate(token.Header["kid"])
		if err != nil {
			panic(err.Error())
		}

		result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			panic(err.Error())
		}

		return result, nil
	}
}
