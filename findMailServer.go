package emailVerifier

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func findMailServer(companyName string) (map[string]*net.MX, error) {
	tlds := []string{".co.uk", ".com", ".net", ".org"} // TODO: Expand this list.
	var ci []string
	var cn []string
	companyWords := strings.Split(companyName, " ")
	for _, w := range companyWords {
		if len(w) > 0 {
			ci = append(ci, fmt.Sprint(w[0]))
		}
		cn = append(cn, w)
	}
	companyInitials := strings.Join(ci, "")
	flatCompanyName := strings.Join(cn, "")
	hosts := map[string]*net.MX{}

	for _, t := range tlds {
		nameHost := strings.Join([]string{flatCompanyName, t}, "")
		initialHost := strings.Join([]string{companyInitials, t}, "")
		res, err := net.LookupMX(nameHost)
		if err == nil {
			hosts[nameHost] = res[0]
		}
		iRes, err := net.LookupMX(initialHost)
		if err == nil {
			hosts[initialHost] = iRes[0]
		}
	}

	if len(hosts) > 0 {
		return hosts, nil
	}

	return map[string]*net.MX{}, errors.New("No MX records for company name or initials: " + companyName)
}
