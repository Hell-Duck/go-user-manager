package main

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func TestDatabaseConnection(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		t.Fatal("DATABASE_URL is empty")
	}
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		t.Fatalf("Invalid open connection to db, error: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatalf("Db not connection, error: %v", err)
	}
}
