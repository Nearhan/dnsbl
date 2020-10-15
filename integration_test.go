// +build integration

package dnsbl

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

const testData = `
INSERT INTO ipdetails (id, created_at, updated_at, response_code, ip_address)
VALUES
  ("4aeec477-2e9d-4cd7-8840-6d3a0a8b7a1b", "2020-10-13T20:56:32+00:00", "2020-10-13T20:56:32+00:00", "3", "1.2.3.4"),
  ("2d3388a2-7106-4116-8ad5-b70dd98b8519", "2020-10-13T20:56:32+00:00", "2020-10-13T20:56:32+00:00", "3", "1.3.3.4"),
  ("d3beeecc-43a5-4891-99b3-bb1aa9d9a0d9", "2020-10-13T20:56:32+00:00", "2020-10-13T20:56:32+00:00", "3", "1.4.3.4")
`

func setupDB(t *testing.T, db *sql.DB) {
	sqls := []string{schema, testData}

	for _, s := range sqls {
		stmt, err := db.Prepare(s)
		if err != nil {
			t.Fatal(err)
		}

		_, err = stmt.Exec()
		if err != nil {
			t.Fatal(err)
		}

		err = stmt.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

}

func tearDownDB() {
	os.Remove("testdata/test.db")
}

func TestGetIPdetail(t *testing.T) {

	db, err := sql.Open("sqlite3", "testdata/test.db")
	require.NoError(t, err)

	//setup db
	setupDB(t, db)
	defer tearDownDB()

	act, err := GetIPDetail(context.Background(), db, "1.2.3.4")
	require.NoError(t, err)

	ti, err := time.Parse(time.RFC3339, "2020-10-13T20:56:32+00:00")
	require.NoError(t, err)

	exp := &IPDetail{
		ID:           "4aeec477-2e9d-4cd7-8840-6d3a0a8b7a1b",
		CreatedAt:    ti,
		UpdatedAt:    ti,
		ResponseCode: "3",
		IPAddress:    "1.2.3.4",
	}

	require.Equal(t, exp, act)
}

func TestListGetIPdetail(t *testing.T) {

	db, err := sql.Open("sqlite3", "testdata/test.db")
	require.NoError(t, err)

	//setup db
	setupDB(t, db)
	defer tearDownDB()

	act, err := listIPDetails(context.Background(), db, []string{"1.2.3.4", "1.3.3.4"})
	require.NoError(t, err)

	ti, err := time.Parse(time.RFC3339, "2020-10-13T20:56:32+00:00")
	require.NoError(t, err)

	exp := []IPDetail{
		{
			ID:           "4aeec477-2e9d-4cd7-8840-6d3a0a8b7a1b",
			CreatedAt:    ti,
			UpdatedAt:    ti,
			ResponseCode: "3",
			IPAddress:    "1.2.3.4",
		},
		{
			ID:           "2d3388a2-7106-4116-8ad5-b70dd98b8519",
			CreatedAt:    ti,
			UpdatedAt:    ti,
			ResponseCode: "3",
			IPAddress:    "1.3.3.4",
		},
	}

	require.Equal(t, exp, act)
}

func TestEnqueue(t *testing.T) {

	db, err := sql.Open("sqlite3", "testdata/test.db")
	require.NoError(t, err)

	//setup db
	setupDB(t, db)
	defer tearDownDB()

	ips := []string{"1.2.3.4", "1.1.1.1"}
	err = Enqueue(context.Background(), db, ips)
	require.NoError(t, err)

	act, err := listIPDetails(context.Background(), db, ips)
	require.NoError(t, err)
	require.Equal(t, 2, len(act))

}
