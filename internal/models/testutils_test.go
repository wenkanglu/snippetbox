package models

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func newTestDB(t *testing.T) *pgxpool.Pool {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}
	testPgConnEnv := os.Getenv("TEST_POSTGRES_CONN")

	pool, err := pgxpool.New(context.Background(), testPgConnEnv)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(context.Background(), string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer pool.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = pool.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return pool
}
