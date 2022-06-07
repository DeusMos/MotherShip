package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.rd.keyww.com/sw/spike/rmd/cmd/clitools"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/session"
)

var Login = &cobra.Command{
	Use:   "login",
	Short: "Login to the fini server with your active directory credentials",
	Long:  `Login uses your active directory credentials to acquire a session token`,
	RunE:  login,
}

func login(c *cobra.Command, args []string) (err error) {
	client := http.Client{}
	endpoint := "https://fini.key.net/api/session"
	username, password, err := clitools.GetCredentials()
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
		err = clitools.NewHelpError(err.Error())
		return
	}
	_, err = session.FromJSON(bodyBytes)
	if err != nil {
		err = clitools.NewHelpError(err.Error())
		return
	}
	return
}
