package emailVerifier

// FindEmail takes a first, last and company name, finds MX and mail server addresses,
// creates possibilities for email addresses and tries them against the server.
func FindEmail(firstName string, lastName string, companyName string) (string, error) {
	test, err := tryEmails(firstName, lastName, companyName)
	if err != nil {
		return "", err
	}
	return test[0], nil
}