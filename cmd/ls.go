package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.rd.keyww.com/sw/spike/rmd/cmd/clitools"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/machine"
)

// LS represents the 'ls' command for introspecting into a backup archive.
var LS = &cobra.Command{
	Use:   "ls",
	Short: "List all the remote machines.",
	Long:  `List provides a grep friendly table of the machines that are available on rmd`,
	RunE:  lsRun,
}

//lsRun prints a table with all the remote machines available on rmd.
func lsRun(cmd *cobra.Command, args []string) (err error) {
	// Validate arguments.
	endpoint := "http://fini.key.net/api/rmd/ls"
	var client = http.Client{}
	l := len(args)
	if l != 0 {
		err = clitools.NewHelpError(
			"'ls' does not take any arguments",
		)
		return
	}
	session, err := ioutil.ReadFile("session")
	if err != nil {
		log.Fatal(err)
	}
	cred := func(r *http.Request) {
		r.Header.Set("Fini-Session", string(session))
	}
	var req *http.Request

	if req, err = http.NewRequest("GET", endpoint, nil); nil != err {
		err = fmt.Errorf(
			"while creating upload request for first chunk got: %v", err,
		)
		return
	}
	cred(req)
	var resp *http.Response
	if resp, err = client.Do(req); nil != err {
		err = fmt.Errorf("while downloading backup got: %v", err)
		return
	}
	defer resp.Body.Close()
	machines := make([]machine.Machine, 0)
	err = json.NewDecoder(resp.Body).Decode(&machines)
	if nil != err {
		fmt.Print(err.Error())
	}
	for _, m := range machines {
		fmt.Println(m)
	}
	return
}
