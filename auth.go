package dnsbl

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// BasicAuthMiddleWear basic auth middle wear
func BasicAuthMiddleWear(next http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !isUserAuthorized(db, user, pass) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"code": "401", "status": "unauthorized"}`)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isUserAuthorized(db *sql.DB, user, pass string) bool {
	stmt := fmt.Sprintf("SELECT pass FROM users WHERE user='%s'", user)
	row, err := db.Query(stmt)
	if err != nil {
		log.Println(err)
	}

	var dbpass string
	var auth bool

	for row.Next() {
		err = row.Scan(&dbpass)
		if err != nil {
			log.Println(err)
		}

		defer row.Close()
		if bcrypt.CompareHashAndPassword([]byte(dbpass), []byte(pass)) == nil {
			auth = true
			break
		}
	}
	return auth
}

func hashAndSalt(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
