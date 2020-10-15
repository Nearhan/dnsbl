package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Nearhan/dnsbl/dnsbl"
)

func main() {

	db, err := sql.Open("sqlite3", "ipdetails.db")
	if err != nil {
		log.Fatal("unable to create database", err)
	}

	if err := dnsbl.InitSchmea(db); err != nil {
		log.Fatal("unable to init database schema", err)
	}

	resolver := dnsbl.NewResolver(db)

	srv := handler.NewDefaultServer(dnsbl.NewExecutableSchema(dnsbl.Config{Resolvers: resolver}))
	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
