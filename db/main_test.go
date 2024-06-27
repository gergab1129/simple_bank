package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	testQueries *Queries
	testStore   *Store
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	conn, err := pgxpool.New(
		ctx,
		"user=root password=jg.5s0o4ht7 host=localhost port=5432 database=simple_bank sslmode=disable",
	)
	if err != nil {
		log.Fatalf("error creating connection to database. got=%s", err)
	}

	testQueries = New(conn)

	testStore = NewStore(conn)

	os.Exit(m.Run())
}
