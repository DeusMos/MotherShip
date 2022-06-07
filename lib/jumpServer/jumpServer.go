package jumpServer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	portforward "gitlab.rd.keyww.com/sw/spike/rmd/lib/portForward"
)

type JumpServer struct {
	Url         string
	JumpSSHPort uint
	SSHPort     uint
	VNCPort     uint
	BrowserPort uint
	ID_RSA      string
	ID_RSA_PUB  string
}

func NewJumpServer(url, iD_RSA, iD_RSA_PUB string, jumpSSHPort, sSHPort, vNCPort, browserPort uint) (js JumpServer) {
	js = JumpServer{Url: url, JumpSSHPort: jumpSSHPort, SSHPort: sSHPort, VNCPort: vNCPort, BrowserPort: browserPort, ID_RSA: iD_RSA, ID_RSA_PUB: iD_RSA_PUB}
	return
}

// Connect will connect all the set ports to the jump server
func (s *JumpServer) Connect() (active bool, err error) {
	jp := fmt.Sprint(s.JumpSSHPort)
	sp := fmt.Sprint(s.SSHPort)
	vp := fmt.Sprint(s.VNCPort)
	bp := fmt.Sprint(s.BrowserPort)

	err = ioutil.WriteFile("/tmp/id_rsa_"+jp, []byte(s.ID_RSA), 0600)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("/tmp/id_rsa_"+jp+".pub", []byte(s.ID_RSA_PUB), 0600)
	if err != nil {
		return
	}
	active = false
	var forwards []portforward.PortForward
	if s.SSHPort > 0 {
		forwards = append(forwards, portforward.PortForward{LocalPort: "22", RemotePort: sp})
		if err != nil {
			return
		}
	}
	if s.VNCPort > 0 {
		forwards = append(forwards, portforward.PortForward{LocalPort: "5910", RemotePort: vp})
		if err != nil {
			return
		}
		active = true
	}
	if s.BrowserPort > 0 {

		forwards = append(forwards, portforward.PortForward{LocalPort: "80", RemotePort: bp})

		if err != nil {
			return
		}
		active = true
	}
	if active {
		connect(jp, s.Url, forwards)
		active = false
	}
	return
}
func connect(jumpPort, url string, forwards []portforward.PortForward) {
	args := []string{"-N", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null"}
	for _, f := range forwards {
		args = append(args, "-R", f.RemotePort+":localhost:"+f.LocalPort)
	}
	args = append(args, "-i", "/tmp/id_rsa_"+jumpPort, "-p", jumpPort, "root@"+url)
	fmt.Println(args)
	ssh_cmd := exec.Command("/usr/bin/ssh", args...)
	ssh_cmd.Stdout = os.Stdout
	ssh_cmd.Stdin = os.Stdin
	ssh_cmd.Stderr = os.Stderr
	err := ssh_cmd.Run()
	if err != nil {
		fmt.Printf("ssh tunnel failed with %+v ", args)
		err = nil
	}

}
