package auth0

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

// TestBuildJwksWellKnownURL verifies that URL to OIDC well-known metadata is correctly built
func TestBuildJwksWellKnownURL(t *testing.T) {
	expectedURL := "https://rainerdemo.eu.auth0.com/.well-known/jwks.json"
	url, err := buildJwksWellKnownURL("https://rainerdemo.eu.auth0.com")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if url != expectedURL {
		t.Errorf("Expected %s, received URL %s", expectedURL, url)
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func newTestClientForURLs(knownURLs map[string]string) *http.Client {
	return newTestClient(func(req *http.Request) *http.Response {
		content, ok := knownURLs[req.URL.String()]
		if !ok {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(content)),
			Header:     make(http.Header),
		}
	})
}

// TestGetPemCertificate verifies that cert is correctly read from OIDC metadata
func TestGetPemCertificate(t *testing.T) {
	const oidcBaseURL = "http://dummy"
	const jwksURL = oidcBaseURL + "/.well-known/jwks.json"
	const certContent = "CERT"
	const kidValue = "KID"
	const jwksResponse = `{ "keys": [ { "kid": "` + kidValue + `", "x5c": [ "` + certContent + `" ] } ] }`
	const expectedCert = "-----BEGIN CERTIFICATE-----\n" + certContent + "\n-----END CERTIFICATE-----"

	pcr := PemCertificateOidcJwksReader{
		HTTPClient:  newTestClientForURLs(map[string]string{jwksURL: jwksResponse}),
		BaseOidcURL: oidcBaseURL,
	}
	cert, err := pcr.GetPemCertificate(kidValue)
	if err != nil {
		t.Errorf("Error getting PEM certificate %v", err)
	}

	if cert != expectedCert {
		t.Errorf("Expected %s, got %s", expectedCert, cert)
	}
}

// TestGetPemCertificateMissingKey checks that not found kid header is correctly handled
func TestGetPemCertificateMissingKey(t *testing.T) {
	const oidcBaseURL = "http://dummy"
	const jwksURL = oidcBaseURL + "/.well-known/jwks.json"
	const jwksResponse = `{ "keys": [ ] }`

	pcr := PemCertificateOidcJwksReader{
		HTTPClient:  newTestClientForURLs(map[string]string{jwksURL: jwksResponse}),
		BaseOidcURL: oidcBaseURL,
	}
	if _, err := pcr.GetPemCertificate("KID"); err == nil {
		t.Error("Expected error did not happen")
	}
}

// TestGetPemCertificateJwksNotFound checks that a 404 when reading OIDC metadata is correctly handled
func TestGetPemCertificateJwksNotFound(t *testing.T) {
	const oidcBaseURL = "http://dummy"

	pcr := PemCertificateOidcJwksReader{
		HTTPClient:  newTestClientForURLs(map[string]string{}),
		BaseOidcURL: oidcBaseURL,
	}
	if _, err := pcr.GetPemCertificate("KID"); err == nil {
		t.Error("Expected error did not happen")
	}
}

// TestValidationKeyGetterAudienceCheck verifies that audience is checked
func TestValidationKeyGetterAudienceCheck(t *testing.T) {
	validationFunc := ValidationKeyGetter("aud", "iss", nil)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"aud": "xyz"})
	if _, err := validationFunc(jwtToken); err == nil {
		t.Error("Expected audience error did not happen")
	}
}

// TestValidationKeyGetterAudienceCheck verifies that issuer is checked
func TestValidationKeyGetterIssuerCheck(t *testing.T) {
	validationFunc := ValidationKeyGetter("aud", "iss", nil)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"aud": "aud", "iss": "xyz"})
	if _, err := validationFunc(jwtToken); err == nil {
		t.Error("Expected audience error did not happen")
	}
}

type PemCertificateErrorReader struct {
}

func (PemCertificateErrorReader) GetPemCertificate(_ interface{}) (string, error) {
	return "", errors.New("Something bad happened")
}

// TestValidationKeyGetterPanicsCertReadError verifies that code panics if cert cannot be read
func TestValidationKeyGetterPanicsCertReadError(t *testing.T) {
	validationFunc := ValidationKeyGetter("aud", "iss", PemCertificateErrorReader{})

	recovered := func() (r bool) {
		r = false
		defer func() {
			if rec := recover(); rec != nil {
				r = true
			}
		}()
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"aud": "aud", "iss": "iss"})
		validationFunc(jwtToken)
		return
	}

	if !recovered() {
		t.Errorf("The code did not panic")
	}
}

type PemCertificateDummyCertReader struct {
	Certificate string
}

func (pcr PemCertificateDummyCertReader) GetPemCertificate(_ interface{}) (string, error) {
	return pcr.Certificate, nil
}

// TestValidationKeyGetterInvalidCert verifies that code panics if cert cannot be parsed
func TestValidationKeyGetterInvalidCert(t *testing.T) {
	validationFunc := ValidationKeyGetter("aud", "iss", PemCertificateDummyCertReader{
		Certificate: `-----BEGIN CERTIFICATE-----\nDUMMY\n-----END CERTIFICATE-----`,
	})

	recovered := func() (r bool) {
		r = false
		defer func() {
			if rec := recover(); rec != nil {
				r = true
			}
		}()
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"aud": "aud", "iss": "iss"})
		validationFunc(jwtToken)
		return
	}

	if !recovered() {
		t.Errorf("The code did not panic")
	}
}

// TestValidationKeyGetter verifies that correct certificate is parsed without error
func TestValidationKeyGetter(t *testing.T) {
	pcr := PemCertificateDummyCertReader{
		Certificate: `-----BEGIN CERTIFICATE-----
MIIDCTCCAfGgAwIBAgIJCCnZ4NypkYWrMA0GCSqGSIb3DQEBCwUAMCIxIDAeBgNVBAMTF3JhaW5lcmRlbW8uZXUuYXV0aDAuY29tMB4XDTE3MDYxODA4MDQ1MVoXDTMxMDIyNTA4MDQ1
MVowIjEgMB4GA1UEAxMXcmFpbmVyZGVtby5ldS5hdXRoMC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC9dO4tyTjdRm5EEsO5iI0OI0jGzBBgwvmsX9vYkybh71Pv
E0KV+/cLhRwPHRmTjMf5pY3iwXGgesbY8YrbA6Q7+jaePBEnlqMZbqSwSFfSv5z/rcGCy3+M+kHkQnNFnha4HFX1+mvTHi+9bsSa+/PO8SVcKBfB1+Ga7OLReNBpVeGYhqVhf0kKy7NG
CTW+0CHQjMgMT7ZASfBquMp7TC/Ol23j3ZFaRbLHd4KKHZDQAJLCG2Sjj//rmIEG1YIRGraUsnoAltjLbTPk+19lEhB0iGPbWN/Y8K+Ks4D8DPIOZB1n7xrcz39ZQE4YKu0nzbPM8Xo9
jFqfKMsNmzMIK3+tAgMBAAGjQjBAMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFOF1NItT6wah1kU09ZdkmOD9jKoOMA4GA1UdDwEB/wQEAwIChDANBgkqhkiG9w0BAQsFAAOCAQEA
Q7+UTgE/z8tjhwCzzljvOL12bD9tO8pIhijUm+gr02ib83/xMoCsPvKp3eKh2FgVd8C/KaYYVhldc6BZqyKoUAB878XcthYzgSAlKeusryCn/e8sXq6+20iOetUbY2XQ7y625Vb5m3Yw
77fxfn8Ooro3bQ7YTpaunSZY3b8hSYGxViP+1PCxmSdr9qjH9CO90QJeKxzXW4z/QapHGoO8fvukRmTuPRfmorqPOlo//Uf+gQMq+qpDYAJNsgcbtlDv3hOTZYRArR8Tl7NTh3GueDHS
5bs5W0eWG4sPkn6TPciISOIJZeAE9J4rplmy1+lNQ947dDg7POh02dX7bZc6lA==
-----END CERTIFICATE-----`,
	}
	validationFunc := ValidationKeyGetter("aud", "iss", pcr)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"aud": "aud", "iss": "iss"})
	if result, err := validationFunc(jwtToken); result == nil || err != nil {
		t.Errorf("Could not parse certificate (%s)", err.Error())
	}
}
