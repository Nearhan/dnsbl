package dnsbl

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Schema for application
var createIpDetailTable = `CREATE TABLE IF NOT EXISTS ipdetails (
    id TEXT PRIMARY KEY,
    created_at TEXT,
    updated_at TEXT,
    response_code TEXT,
    ip_address TEXT
)`

var createUsersTable = `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	user TEXT,
	pass TEXT
)`

var insertBaseUser = `INSERT INTO users (user, pass) VALUES ('%s', '%s');
`

//InitSchema inits database schema
func InitSchema(db *sql.DB) error {

	sqls := []string{
		createIpDetailTable,
		createUsersTable,
		fmt.Sprintf(insertBaseUser, "secureworks", hashAndSalt([]byte("password"))),
	}

	for _, s := range sqls {
		stmt, err := db.Prepare(s)
		if err != nil {
			return err
		}

		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func listIPQuery(ips []string) string {

	args := []string{}
	q := "SELECT * FROM ipdetails WHERE ip_address IN (%s)"
	for _, ip := range ips {
		args = append(args, fmt.Sprintf("'%s'", ip))
	}

	return fmt.Sprintf(q, strings.Join(args, ","))
}

func makeInsertStmt(ipds []IPDetail) string {

	values := []string{}
	for _, ipd := range ipds {
		values = append(values, ipd.toInsertSQL())
	}

	q := "INSERT INTO ipdetails (id, created_at, updated_at, response_code, ip_address) VALUES %s"
	return fmt.Sprintf(q, strings.Join(values, ","))

}

func makeUpdateStmt(ipds []IPDetail) string {
	stmts := []string{}
	for _, ipd := range ipds {
		stmts = append(stmts, fmt.Sprintf(
			"Update ipdetails SET updated_at='%s', response_code='%s' WHERE ip_address='%s'",
			ipd.UpdatedAt.UTC().Format(time.RFC3339),
			ipd.ResponseCode,
			ipd.IPAddress,
		))
	}
	return strings.Join(stmts, ",")
}

// IPDetail holds information about an ip address lookup
type IPDetail struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ResponseCode string    `json:"response_code"`
	IPAddress    string    `json:"ip_address"`
}

func (i IPDetail) toInsertSQL() string {

	return fmt.Sprintf(
		"('%s', '%s', '%s', '%s', '%s')",
		i.ID,
		i.CreatedAt.UTC().Format(time.RFC3339),
		i.UpdatedAt.UTC().Format(time.RFC3339),
		i.ResponseCode,
		i.IPAddress,
	)
}

//diffIpDetails separates into insert and update collections
func diffIPDetails(newDetails, foundDetails []IPDetail) ([]IPDetail, []IPDetail) {
	insert := []IPDetail{}
	update := []IPDetail{}
	ndm := map[string]IPDetail{}
	for _, nd := range newDetails {
		ndm[nd.IPAddress] = nd
	}

	for _, fd := range foundDetails {
		if ipd, ok := ndm[fd.IPAddress]; ok {
			fd.UpdatedAt = ipd.UpdatedAt
			fd.ResponseCode = ipd.ResponseCode
			delete(ndm, fd.IPAddress)
			update = append(update, fd)
		}
	}

	for _, ipd := range ndm {
		insert = append(insert, ipd)
	}

	return insert, update
}
