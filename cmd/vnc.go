package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.rd.keyww.com/sw/spike/rmd/cmd/clitools"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/fini"
)

// LS represents the 'ls' command for introspecting into a backup archive.
var VNC = &cobra.Command{
	Use:   "vnc",
	Short: "Builds an vnc tunnel to a jump server, then waits for the target machine to connect",
	Long:  `vnc takes a machine id and then tells fini to have it connect to a jump server where we will be wating.`,
	RunE:  VNCRun,
}

// VNC sets us up with a vnc connection to target machine.
func VNCRun(cmd *cobra.Command, args []string) (err error) {
	endpoint := "http://fini.key.net/api/rmd/vnc"
	js, err := fini.GetJumpServer(endpoint, string(args[0]))
	if err != nil {
		err = clitools.NewHelpError(err.Error())
		return
	}
	p := fmt.Sprint(js.SSHPort)
	err = ioutil.WriteFile("./id_rsa_"+p, []byte(js.ID_RSA), 0644)
	if err != nil {
		err = clitools.NewHelpError(err.Error())
		return
	}
	vp := fmt.Sprint(js.VNCPort)
	vnc_cmd := exec.Command("vnc", "-g", "-L", vp+":localhost:"+vp, "-i", "./id_rsa"+p, "-p", p, "user@"+js.Url)
	vnc_cmd.Stdout = os.Stdout
	vnc_cmd.Stdin = os.Stdin
	vnc_cmd.Stderr = os.Stderr
	vnc_cmd.Run()
	return
}
