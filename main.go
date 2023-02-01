package main

import (
	"fmt"
	"context"
	"net/http"

	// "flag"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rail44/g/openapi"
	"github.com/rail44/g/sqlc/generated"
)

type Secret struct {
	SupabaseKey string `toml:"supabaseKey"`
}

type GraphqlTranstport struct {
	supabaseKey *string
}

func (trastport *GraphqlTranstport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("apikey", *trastport.supabaseKey)
	return http.DefaultTransport.RoundTrip(req)
}

type Server struct {
        queries *sqlc.Queries
}

func (server Server) GetAccountAccountIdBalance(ctx context.Context, req openapi.GetAccountAccountIdBalanceRequestObject) (openapi.GetAccountAccountIdBalanceResponseObject, error) {
        account, err := server.queries.GetAccount(ctx, 1)
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println(account)
	balance := "100"
	res := openapi.GetAccountAccountIdBalance200JSONResponse{
		Balance: &balance,
	}
	return res, nil
}

func (server Server) PostAccountRegister(ctx context.Context, req openapi.PostAccountRegisterRequestObject) (openapi.PostAccountRegisterResponseObject, error) {
	accountId := 1
	res := openapi.PostAccountRegister200JSONResponse{
		AccountId: &accountId,
	}
	return res, nil
}

func main() {
        db, err := sql.Open("postgres", "user=postgres dbname=g password=password host=localhost sslmode=disable")
        if err != nil {
                panic(err)
        }
        queries := sqlc.New(db)

        server := Server{ queries: queries }
	http.ListenAndServe(":3000", openapi.Handler(openapi.NewStrictHandler(server, nil)))
}
