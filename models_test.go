package dnsbl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestListIPQuery(t *testing.T) {

	ips := []string{"1.2.3.4", "1.2.3.5"}
	act := listIPQuery(ips)
	exp := "SELECT * FROM ipdetails WHERE ip_address IN ('1.2.3.4','1.2.3.5')"
	require.Equal(t, exp, act)
}

func TestMakeInsertStmt(t *testing.T) {

	ti, err := time.Parse(time.RFC3339, "2020-10-14T15:00:21Z")
	require.NoError(t, err)
	ipds := []IPDetail{
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

	act := makeInsertStmt(ipds)
	exp := `INSERT INTO ipdetails (id, created_at, updated_at, response_code, ip_address) VALUES ('4aeec477-2e9d-4cd7-8840-6d3a0a8b7a1b', '2020-10-14T15:00:21Z', '2020-10-14T15:00:21Z', '3', '1.2.3.4'),('2d3388a2-7106-4116-8ad5-b70dd98b8519', '2020-10-14T15:00:21Z', '2020-10-14T15:00:21Z', '3', '1.3.3.4')`
	require.Equal(t, exp, act)
}

func TestMakeUpdateStmt(t *testing.T) {

	ti, err := time.Parse(time.RFC3339, "2020-10-14T15:00:21Z")
	require.NoError(t, err)
	ipds := []IPDetail{
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

	act := makeUpdateStmt(ipds)
	exp := "Update ipdetails SET updated_at='2020-10-14T15:00:21Z', response_code='3' WHERE ip_address='1.2.3.4',Update ipdetails SET updated_at='2020-10-14T15:00:21Z', response_code='3' WHERE ip_address='1.3.3.4'"
	require.Equal(t, exp, act)
}

func TestDiffIPDetail(t *testing.T) {
	ti, err := time.Parse(time.RFC3339, "2020-10-14T15:00:21Z")
	require.NoError(t, err)
	newIPd := []IPDetail{
		{
			ID:           "ff9ec662-cfe9-49ae-9372-29e711821fa9",
			CreatedAt:    ti.Add(1 * time.Hour),
			UpdatedAt:    ti.Add(1 * time.Hour),
			ResponseCode: "2",
			IPAddress:    "1.2.3.4",
		},
	}

	foundIPd := []IPDetail{
		{
			ID:           "4aeec477-2e9d-4cd7-8840-6d3a0a8b7a1b",
			CreatedAt:    ti,
			UpdatedAt:    ti,
			ResponseCode: "3",
			IPAddress:    "1.2.3.4",
		},
	}

	mergedUpdate := IPDetail{
		ID:           "4aeec477-2e9d-4cd7-8840-6d3a0a8b7a1b",
		CreatedAt:    ti,
		UpdatedAt:    ti.Add(1 * time.Hour),
		ResponseCode: "2",
		IPAddress:    "1.2.3.4",
	}

	insert, update := diffIPDetails(newIPd, foundIPd)
	require.Empty(t, insert)
	require.Len(t, update, 1)
	require.Equal(t, mergedUpdate, update[0])

}
