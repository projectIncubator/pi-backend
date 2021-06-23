package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func InitAuthMiddleware() *jwtmiddleware.JWTMiddleware {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			url := os.Getenv("ISS")
			// Verify 'aud' claim
			aud := url + "api/v2/"
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			iss := url
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	return jwtMiddleware
}
func getPemCert(token *jwt.Token) (string, error) {
	url := os.Getenv("ISS")
	cert := ""
	resp, err := http.Get(url + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

type Scope int

const (
	USER Scope = 1 + iota
	ADMIN
	CREATOR
)

type AuthWraper struct {
	id string
}

func (app *App) middleware(next http.Handler, scope Scope) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userToken := r.Header.Get("User-ID")

		switch scope {
		case USER:
			usr, err := app.store.ScopeProvider.GetUserID(userToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			if !usr.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			auth := AuthWraper{usr.String}
			ctx := context.WithValue(r.Context(), "user_id", auth)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		case ADMIN:
			projectID := mux.Vars(r)["proj_id"]
			usr, err := app.store.ScopeProvider.GetAdminID(userToken, projectID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			if !usr.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			auth := AuthWraper{usr.String}
			ctx := context.WithValue(r.Context(), "user_id", auth)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		case CREATOR:
			projectID := mux.Vars(r)["proj_id"]
			usr, err := app.store.ScopeProvider.GetCreatorID(userToken, projectID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			if !usr.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}

			auth := AuthWraper{usr.String}
			ctx := context.WithValue(r.Context(), "user_id", auth)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	})
}
