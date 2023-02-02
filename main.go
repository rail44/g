package main

import (
	"net/http"

	// "flag"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rail44/g/accounts"
        "github.com/go-chi/chi/v5"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=g password=password host=localhost sslmode=disable")
	if err != nil {
		panic(err)
	}

        r := chi.NewRouter()
	accountsCotroller := accounts.NewController(db)
	r.Mount("/accounts", accountsCotroller)
	http.ListenAndServe(":3000", r)
}
