package main

import (
	"context"
	"fmt"
	"homework-7/internal/pkg/db"
	"homework-7/internal/pkg/repository/postgresql"
	"homework-7/internal/pkg/server"
	"log"
	"net/http"
)

const port = ":9000"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connection string
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DBname)

	database, err := db.NewDB(ctx, psqlConn)
	defer database.GetPool(ctx).Close()

	userRepo := postgresql.NewUsers(database)

	if err != nil {
		return
	}

	mux := server.CreateUserServer(ctx, userRepo)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
