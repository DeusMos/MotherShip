package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gitlab.rd.keyww.com/sw/spike/rmd/cmd/clitools"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/fini"
)

// LS represents the 'ls' command for introspecting into a backup archive.
var Proxy = &cobra.Command{
	Use:   "proxy",
	Short: "Build an ssh tunnel to a jump server with http forwared, then waits for the target machine to connect",
	Long:  `ssh takes a machine id and then tells fini to have it connect to a jump server where we will be wating.`,
	RunE:  ProxyRun,
}

// ProxyRun sets us up with a ssh connection to target machine.
func ProxyRun(cmd *cobra.Command, args []string) (err error) {
	endpoint := "https://fini.key.net/api/rmd/proxy"
	js, err := fini.GetJumpServer(endpoint, string(args[0]))

	p := fmt.Sprint(js.SSHPort)
	if err != nil {
		err = clitools.NewHelpError(err.Error())
		return
	}

	vnc_cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-D", p+":localhost:"+p, "-i", "./id_rsa"+p, "-p", p, "user@"+js.Url)
	vnc_cmd.Stdout = os.Stdout
	vnc_cmd.Stdin = os.Stdin
	vnc_cmd.Stderr = os.Stderr
	vnc_cmd.Run()
	return
}
