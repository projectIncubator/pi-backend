package routes

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Scope int

const (
	USER Scope = 1 + iota
	ADMIN
	CREATOR
)

type AuthWraper struct {
	id string;
};



func (app *App) middleware(next http.Handler, scope Scope) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			userToken := authHeader[1]

			switch(scope){
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
		}

		return
	})
}