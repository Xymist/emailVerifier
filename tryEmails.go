package emailVerifier

import "net/textproto"
import "errors"

import "strings"

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

func setupMX(conn *textproto.Conn, fromEmail string) error {
	if err := checkResponse(conn, "", 220); err != nil {
		return errors.New("Could not establish connection: " + err.Error())
	}

	if err := checkResponse(conn, "HELO HI", 250); err != nil {
		return errors.New("Did not receive HELO response: " + err.Error())
	}

	if err := checkResponse(conn, "mail from: <"+fromEmail+">", 250); err != nil {
		return errors.New("Mail from " + fromEmail + " not accepted: " + err.Error())
	}
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

	if err := setupMX(conn, potentialEmails[0]); err != nil {
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
