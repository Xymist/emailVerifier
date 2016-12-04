package emailVerifier

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func tryEmails(firstName string, lastName string, companyName string) ([]string, error) {
	potentialEmails := []string{}
	mx, err := findMailServer(companyName)
	if err != nil {
		return []string{}, err
	}
	for s := range mx {
		for _, email := range generateEmails(firstName, lastName, s) {
			potentialEmails = append(potentialEmails, email)
		}
	}

	foundEmails := []string{}
	for _, e := range potentialEmails {
		//TODO: This step is slow, switch to using goroutines and a channel.
		if err := VerifyEmail(e); err != nil {
			continue
		}
		foundEmails = append(foundEmails, e)
	}

	if len(foundEmails) > 0 {
		return foundEmails, nil
	}
	return []string{}, errors.New("Nothing Found!")
}

func generateEmails(firstName string, lastName string, mailserver string) []string {
	firstInitial := fmt.Sprint(string(firstName[0]))
	emailAddresses := []string{
		firstName + "." + lastName + "@" + mailserver,
		firstInitial + lastName + "@" + mailserver,
		lastName + firstInitial + "@" + mailserver,
		firstName + lastName + "@" + mailserver,
		firstName + "@" + mailserver}

	return emailAddresses
}

func findMailServer(companyName string) (map[string]*net.MX, error) {
	tlds := []string{".co.uk", ".com", ".net", ".org", ".io"} // TODO: Expand this list.
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
