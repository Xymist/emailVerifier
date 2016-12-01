package emailVerifier

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func findMailServer(companyName string) (string, *net.MX, error) {
	tlds := []string{".co.uk", ".com", ".net", ".org"} // TODO: Expand this list.
	var ci []string
	var cn []string
	var err error
	companyWords := strings.Split(companyName, " ")
	for _, w := range companyWords {
		ci = append(ci, fmt.Sprint(w[0]))
		cn = append(cn, w)
	}
	companyInitials := strings.Join(ci, "")
	flatCompanyName := strings.Join(cn, "")

	for _, t := range tlds {
		nameHost := strings.Join([]string{flatCompanyName, t}, "")
		initialHost := strings.Join([]string{companyInitials, t}, "")
		res, err := net.LookupMX(nameHost)
		if err == nil {
			return nameHost, res[0], nil
		}
		iRes, err := net.LookupMX(initialHost)
		if err == nil {
			return initialHost, iRes[0], nil
		}
	}
	return "", nil, errors.New(err.Error() + ": No MX records for company name or initials")
}
