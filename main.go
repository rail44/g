package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/rail44/g/accounts"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=g password=password host=localhost sslmode=disable")
	if err != nil {
		log.Fatal(fmt.Sprintf("open postgres: %w", err))
	}

	r := chi.NewRouter()
	accountsCotroller := accounts.NewController(db)
	r.Mount("/accounts", accountsCotroller)

	log.Print("start listening")
	err = http.ListenAndServe(":3000", r)
	log.Fatal(fmt.Sprintf("listening: %w", err))
}
