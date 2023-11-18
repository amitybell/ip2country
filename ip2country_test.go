package ip2country

import (
	"testing"
)

func TestLookup(t *testing.T) {
	cases := map[string]string{
		"103.2.196.115":   "AU",
		"107.155.124.170": "US",
		"118.195.165.100": "CN",
		"156.146.56.38":   "SG",
		"185.70.107.18":   "RU",
		"5.188.120.51":    "ZA",
		"51.89.235.89":    "GB",
	}
	for ip, exp := range cases {
		got, err := LookupString(ip)
		if err != nil {
			t.Fatalf("%s: %s", ip, err)
		}
		if got != exp {
			t.Fatalf("%s: Expected %s; Got %s", ip, exp, got)
		}
	}
}
