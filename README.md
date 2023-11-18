# Intro

Package `ip2country` provides IPv4 and IPv6 address to country lookups in Go.

# Install

    go get github.com/amitybell/ip2country

# Usage

    package main

    import (
    	"github.com/amitybell/ip2country"
    	"fmt"
    )

    func main() {
    	fmt.Println(ip2country.LookupString("1.2.3.4"))
    }

# License

The package itself is distributed under the MIT license.

The database is compiled from the MaxMind geolite2-city CSVs distributed by https://github.com/sapics/ip-location-db/tree/main/geolite2-country and is therefore subject to the MaxMind EULA https://www.maxmind.com/en/geolite2/eula
