package ip2country

import (
	_ "embed"
	"github.com/oschwald/maxminddb-golang"
	"net"
)

var (
	//go:embed geolite2-country.mmdb
	mmdbSrc []byte

	mmdb = func() *maxminddb.Reader {
		db, err := maxminddb.FromBytes(mmdbSrc)
		if err != nil {
			panic(err)
		}
		return db
	}()
)

// Lookup returns the country (code) in which ip is located, or an error if not found.
func Lookup(ip net.IP) (string, error) {
	var record string
	err := mmdb.Lookup(ip, &record)
	return record, err
}

// LookupString is a wrapper Lookup(net.ParseIP(ip))
func LookupString(ip string) (string, error) {
	return Lookup(net.ParseIP(ip))
}
