package clitools

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/term"
)

// GetCredentials from the user through the command line.
func GetCredentials() (username, password string, err error) {
	fmt.Println("This operation requires authentication for the fini.key.net service.")
	fmt.Println("Please supply your Key Active Directory credentials.")

	// Get username.
	fmt.Print("username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username = scanner.Text()

	// Get password.
	fmt.Print("password: ")
	var bpw []byte
	if bpw, err = term.ReadPassword(int(syscall.Stdin)); nil != err {
		return
	}
	fmt.Println("")
	password = string(bpw)

	return
}

// RegisterMachine will get a session for the machine to post stats to fini.
func RegisterMachine(uniqueName string) (session string, err error) {
	client := http.Client{}
	endpoint := "https://fini.key.net/api/session/machine"
	username, password, err := GetCredentials()
	if err != nil {
		fmt.Println(err)
	}
	cred := func(r *http.Request) {
		r.SetBasicAuth(username, password)
	}
	var req *http.Request
	if req, err = http.NewRequest("GET", endpoint, nil); nil != err {
		err = fmt.Errorf("while registering machine got: %v", err)
		return
	}
	req.Header.Set("Fini-Machine-Name", uniqueName)
	req.Header.Set("Fini-Machine-Type", "g6")
	cred(req)
	var resp *http.Response
	if resp, err = client.Do(req); nil != err {
		err = fmt.Errorf("while downloading session got: %v", err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	session = string(bodyBytes)
	err = ioutil.WriteFile("./session", []byte(session), 0644)
	if err != nil {
		return
	}
	return

}
