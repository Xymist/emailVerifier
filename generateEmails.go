package emailVerifier

import "fmt"

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
