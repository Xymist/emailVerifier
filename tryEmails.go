package emailVerifier

import "net/textproto"
import "errors"

import "strings"
import "net"

func checkResponse(conn *textproto.Conn, request string, code int) error {
	if request != "" {
		req, err := conn.Cmd(request)
		if err != nil {
			return errors.New("Command not accepted: " + request)
		}
		conn.StartResponse(req)
		defer conn.EndResponse(req)
	}
	_, _, err := conn.ReadCodeLine(code)
	if err != nil {
		return errors.New("Did not get intended response (" + string(code) + "): " + err.Error())
	}
	return nil
}

func verifyEmail(email string) error {
	host := strings.Split(email, "@")[1]

	res, err := net.LookupMX(host)
	if err != nil {
		return errors.New("Incorrect Host Address")
	}
	mxServer := res[0]

	return nil
}

func tryEmails(firstName string, lastName string, companyName string) ([]string, error) {
	ms, mx, err := findMailServer(companyName)
	if err != nil {
		return []string{}, err
	}

	potentialEmails := generateEmails(firstName, lastName, ms)
	foundEmails := []string{}
	mxHost := strings.TrimRight(mx.Host, ".")

	conn, err := textproto.Dial("tcp", mxHost+":25")
	if err != nil {
		return []string{}, err
	}

	defer conn.Close()

	if err = checkResponse(conn, "", 220); err != nil {
		return []string{}, err
	}

	if err := checkResponse(conn, "HELO HI", 250); err != nil {
		return []string{}, err
	}

	if err := checkResponse(conn, "mail from: <"+potentialEmails[0]+">", 250); err != nil {
		return []string{}, err
	}

	for _, e := range potentialEmails {

		if err := checkResponse(conn, "rcpt to: <"+e+">", 250); err != nil {
			continue
		}
		foundEmails = append(foundEmails, e)
	}

	return foundEmails, nil
}
