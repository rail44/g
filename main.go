package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"github.com/rail44/g/accounts"
)

type DBConfig struct {
	User *string
	Pass *string
	Host *string
	Port *int
	Name *string
}

func (config DBConfig) postgresUri() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		*config.User,
		*config.Pass,
		*config.Host,
		*config.Port,
		*config.Name,
	)
}

func main() {
	port := flag.Int("port", 0, "Port for g daemon")
	dbConfig := DBConfig{
		Host: flag.String("dbhost", "", "Hostname for postgresql"),
		Port: flag.Int("dbport", 0, "Port number for postgresql"),
		Name: flag.String("dbname", "", "Database name for postgresql"),
		User: flag.String("dbuser", "", "Database name for postgresql"),
		Pass: flag.String("dbpass", "", "Password for postgresql"),
	}
	flag.Parse()

	db, err := sql.Open("postgres", dbConfig.postgresUri())
	if err != nil {
		log.Fatal(fmt.Sprintf("open postgres: %w", err))
	}

	r := chi.NewRouter()
	accountsCotroller := accounts.NewController(accounts.NewModel(db))
	r.Mount("/accounts", accountsCotroller)

	listenAddr := fmt.Sprintf(":%d", *port)

	log.Printf("Listening on %s", listenAddr)
	err = http.ListenAndServe(listenAddr, r)
	log.Fatal(fmt.Sprintf("listening: %w", err))
}
