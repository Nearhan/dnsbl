package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Nearhan/dnsbl"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	db, err := sql.Open("sqlite3", "ipdetails.db")
	if err != nil {
		log.Fatal("unable to create database", err)
	}

	if err := dnsbl.InitSchema(db); err != nil {
		log.Fatal("unable to init database schema", err)
	}

	resolver := dnsbl.NewResolver(db)

	srv := handler.NewDefaultServer(dnsbl.NewExecutableSchema(dnsbl.Config{Resolvers: resolver}))
	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", dnsbl.BasicAuthMiddleWear(srv, db))
	log.Printf("starting graphql server on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
