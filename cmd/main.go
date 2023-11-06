package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"time"
)

type User struct {
	UserId string `json:"userId"`
}

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil) // replace with secret key
}
func main() {
	addr := ":3333"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, router())
}

func router() http.Handler {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Use(jwtauth.Authenticator)

		r.Get("/admin", Protected)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", welcome)
		r.Post("/sin-up", SinUp)
	})

	return r
}
func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome anonymous"))
}

func SinUp(w http.ResponseWriter, r *http.Request) {
	var User User
	json.NewDecoder(r.Body).Decode(&User)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"user_id": User.UserId,
		"exp":     jwtauth.ExpireIn(3 * time.Minute),
	})
	fmt.Fprint(w, tokenString)
}

func Protected(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
}
