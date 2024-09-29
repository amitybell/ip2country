//go:build generate

//go:generate go run gen.go
//go:generate go test

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"go4.org/netipx"
)

func insert(tree *mmdbwriter.Tree, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http.Get: %s: %w", url, err)
	}
	defer resp.Body.Close()

	cr := csv.NewReader(resp.Body)

	for {
		row, err := cr.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("csv.Read: %w", err)
		}

		if len(row) != 3 {
			return fmt.Errorf("Expected 3 columns, got %d: %v", len(row), row)
		}

		from, err := netip.ParseAddr(row[0])
		if err != nil {
			return fmt.Errorf("ParseAddr: %s: %w", row[0], err)
		}

		to, err := netip.ParseAddr(row[1])
		if err != nil {
			return fmt.Errorf("ParseAddr: %s: %w", row[1], err)
		}

		record := mmdbtype.String(row[2])
		for _, pfx := range netipx.IPRangeFrom(from, to).Prefixes() {
			ipnet := netipx.PrefixIPNet(pfx)
			if err := tree.Insert(ipnet, record); err != nil {
				return fmt.Errorf("Insert: %s: %w", ipnet, err)
			}
		}
	}
}

func main() {
	tree := check2(mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType:            "ip2country",
			RecordSize:              24,
			IncludeReservedNetworks: true,
		},
	))

	urls := []string{
		"https://github.com/sapics/ip-location-db/raw/main/geolite2-country/geolite2-country-ipv4.csv",
		"https://github.com/sapics/ip-location-db/raw/main/geolite2-country/geolite2-country-ipv6.csv",
	}
	for _, url := range urls {
		check(insert(tree, url))
	}

	f := check2(os.Create("geolite2-country.mmdb.zst"))
	defer checkClose(f)

	w := check2(zstd.NewWriter(f, zstd.WithEncoderLevel(zstd.SpeedBestCompression)))
	defer checkClose(w)

	check2(tree.WriteTo(w))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkClose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

func check2[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
