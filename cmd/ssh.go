package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"gitlab.rd.keyww.com/sw/spike/rmd/cmd/clitools"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/fini"
)

// LS represents the 'ls' command for introspecting into a backup archive.
var SSH = &cobra.Command{
	Use:   "ssh",
	Short: "Builds an ssh tunnel to a jump server, then waits for the target machine to connect",
	Long:  `ssh takes a machine id and then tells fini to have it connect to a jump server where we will be wating.`,
	RunE:  SSHRun,
}

// SSHRun sets us up with a ssh connection to target machine.
func SSHRun(cmd *cobra.Command, args []string) (err error) {

	fmt.Println(args)
	endpoint := "http://woodfam.us:8888/api/ssh"
	l := len(args)
	fmt.Println(l)
	if l != 1 {
		err = clitools.NewHelpError(
			"'ssh' takes only one argument, the machine id.",
		)
		return
	}

	js, err := fini.GetJumpServer(endpoint, string(args[0]))
	if err != nil {
		err = clitools.NewHelpError(fmt.Errorf("while connecting got: %v", err).Error())
		return
	}
	jp := fmt.Sprint(js.JumpSSHPort)
	sp := fmt.Sprint(js.SSHPort)
	err = ioutil.WriteFile("/tmp/id_rsa_"+jp, []byte(js.ID_RSA), 0600)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("/tmp/id_rsa_"+jp+".pub", []byte(js.ID_RSA_PUB), 0600)
	if err != nil {
		return
	}
	//todo make this suck less
	time.Sleep(2 * time.Second)

	//check if the localhost has freeports don't assume they are all free.
	fmt.Println("ssh", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-L", sp+":localhost:"+sp, "-i", "/tmp/id_rsa_"+jp, "-p", jp, "root@"+js.Url)
	vnc_cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-L", sp+":localhost:"+sp, "-i", "/tmp/id_rsa_"+jp, "-p", jp, "root@"+js.Url)
	vnc_cmd.Stdout = os.Stdout
	vnc_cmd.Stdin = os.Stdin
	vnc_cmd.Stderr = os.Stderr
	vnc_cmd.Run()
	return
}
