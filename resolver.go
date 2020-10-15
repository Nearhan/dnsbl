package dnsbl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

//ErrInvalidIP not a valid ip error
var ErrInvalidIP = errors.New("not valid ip")

//ErrNotFound ip detail not found
var ErrNotFound = errors.New("ip detail not found")

// Resolver holds deps
type Resolver struct {
	db *sql.DB
}

// New
func NewResolver(db *sql.DB) *Resolver {
	return &Resolver{
		db: db,
	}
}

// GetIPDetail function that returns an ip
func GetIPDetail(ctx context.Context, db *sql.DB, ip string) (*IPDetail, error) {

	if !isValidIP4(ip) {
		return nil, ErrInvalidIP
	}

	ipds, err := listIPDetails(ctx, db, []string{ip})
	if err != nil {
		return nil, err
	}

	if len(ipds) <= 0 {
		return nil, ErrNotFound

	}

	return &ipds[0], nil
}

// GetIPDetail function that returns an ip
func listIPDetails(ctx context.Context, db *sql.DB, ips []string) ([]IPDetail, error) {

	row, err := db.Query(listIPQuery(ips))
	if err != nil {
		return nil, err
	}

	defer row.Close()
	ipds := []IPDetail{}
	for row.Next() {
		var created string
		var updated string
		ipd := IPDetail{}

		err := row.Scan(&ipd.ID, &created, &updated, &ipd.ResponseCode, &ipd.IPAddress)
		if err != nil {
			return nil, err
		}

		ipd.CreatedAt, err = time.Parse(time.RFC3339, created)
		if err != nil {
			return nil, fmt.Errorf("unable to parse created_at field : %v", err)
		}

		ipd.UpdatedAt, err = time.Parse(time.RFC3339, updated)
		if err != nil {
			return nil, fmt.Errorf("unable to parse updated_at field : %v", err)
		}

		ipds = append(ipds, ipd)
	}

	return ipds, nil
}

// Enqueue ....
func Enqueue(ctx context.Context, db *sql.DB, ips []string) error {

	for _, ip := range ips {
		if !isValidIP4(ip) {
			return ErrInvalidIP
		}
	}

	wg := &sync.WaitGroup{}
	newDetails := NewIPDetailSlice()
	for _, ip := range ips {
		wg.Add(1)
		go lookupIP(wg, ip, newDetails)
	}
	wg.Wait()

	foundDetails, err := listIPDetails(ctx, db, ips)
	if err != nil {
		return err
	}

	// diff into update and insert ip details
	insert, update := diffIPDetails(newDetails.ipds, foundDetails)

	_, err = db.Exec(makeInsertStmt(insert))
	if err != nil {
		return fmt.Errorf("error when inserting ip details %s", err)

	}

	_, err = db.Exec(makeUpdateStmt(update))
	if err != nil {
		return fmt.Errorf("error when update ip details %s", err)
	}

	return nil
}

// reverseIP reverses an ip
func reverseIP(ip string) string {
	splitIP := strings.Split(ip, ".")

	for i, j := 0, len(splitIP)-1; i < len(splitIP)/2; i, j = i+1, j-1 {
		splitIP[i], splitIP[j] = splitIP[j], splitIP[i]
	}

	return strings.Join(splitIP, ".")
}

// isValidIp checks if its a valid ip
func isValidIP4(addr string) bool {
	ip := net.ParseIP(addr)
	if ip != nil && ip.To4() != nil {
		return true
	}
	return false
}

func lookupIP(wg *sync.WaitGroup, ip string, ipds *IPDetailSlice) {
	defer wg.Done()
	ti := time.Now().UTC()

	ipd := IPDetail{
		ID:           uuid.NewV4().String(),
		IPAddress:    ip,
		CreatedAt:    ti,
		UpdatedAt:    ti,
		ResponseCode: "",
	}

	ip = reverseIP(ip)
	ip = fmt.Sprintf("%s.%s", ip, "zen.spamhaus.org")
	res, err := net.LookupIP(ip)
	if err != nil {
		log.Printf("error when looking up %s \n %s", ip, err)
	}

	if len(res) > 0 {
		ipd.ResponseCode = res[0].String()
	}

	ipds.Lock()
	ipds.ipds = append(ipds.ipds, ipd)
	defer ipds.Unlock()

}

//IPDetailSlice slice with a lock
type IPDetailSlice struct {
	*sync.Mutex
	ipds []IPDetail
}

//NewIPDetailSlice creates a new IPDetailSlice
func NewIPDetailSlice() *IPDetailSlice {
	return &IPDetailSlice{&sync.Mutex{}, []IPDetail{}}
}
