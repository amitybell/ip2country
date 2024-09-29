package ip2country

import (
	"embed"
	"io"
	"net"
	"sync"

	"github.com/klauspost/compress/zstd"
	"github.com/oschwald/maxminddb-golang"
)

var (
	//go:embed geolite2-country.mmdb.zst
	files embed.FS

	loadDB = sync.OnceValues(initDB)
)

// Lookup returns the country (code) in which ip is located, or an error if not found.
func Lookup(ip net.IP) (string, error) {
	db, err := loadDB()
	if err != nil {
		return "", err
	}

	var record string
	err = db.Lookup(ip, &record)
	return record, err
}

// LookupString is a wrapper Lookup(net.ParseIP(ip))
func LookupString(ip string) (string, error) {
	return Lookup(net.ParseIP(ip))
}

func initDB() (*maxminddb.Reader, error) {
	f, err := files.Open("geolite2-country.mmdb.zst")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := zstd.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return maxminddb.FromBytes(src)
}
