package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/jumpServer"
)

const (
	// docker is the location of the Docker executable.
	docker = "/usr/bin/docker"
)

type containerNS struct{}
type Cntr struct {
	ID string `json:"ID"`
}

var Container *containerNS

// List running Docker Containers.
func (*containerNS) List() (containers []*Cntr, err error) {
	// Run 'docker container ls' so that each entry is printed in standalone
	// JSON.
	format := "{{json .}}"
	cmd := exec.Command(docker, "container", "ls", "--format", format)
	var out []byte
	if out, err = cmd.Output(); nil != err {
		err = fmt.Errorf(
			"while running 'docker container ls': %s | %v", string(out), err,
		)
		return
	}

	// Each line is a JSON entry. Split on newlines.
	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		// Remove any empty lines (one trailing newline).
		if len(line) == 0 {
			continue
		}

		// Parse each JSON entry into a container.
		cntr := &Cntr{}
		if err = json.Unmarshal(line, cntr); nil != err {
			return
		}

		// Add the entry to the list.
		containers = append(containers, cntr)
	}

	return
}

// Prune stopped Docker containers.
func (*containerNS) Prune() (err error) {
	// Run 'docker container prune' to remove dangling Docker Containers.
	cmd := exec.Command(docker, "container", "prune", "--force")
	var out []byte
	if out, err = cmd.CombinedOutput(); nil != err {
		err = fmt.Errorf(
			"while running 'docker container prune': %s | %v", string(out), err,
		)
		return
	}

	return
}

// Wipe all running containers. This is proving necessary in case a container
// refuses to shutdown and the ports need freed.
func (*containerNS) Wipe() (err error) {
	// Get a list of running containers.
	var containers []*Cntr
	if containers, err = Container.List(); nil != err {
		return
	}

	// If there are no containers to wipe, we are finished.
	if len(containers) == 0 {
		return
	}

	// Prepare 'docker container stop' command to allow 10 seconds for a
	// container to stop before being forcefully killed. Give Docker all
	// container IDs simultaneously so that Docker handles the waits in parallel
	// rather than taking up to 60 seconds to stop 6 containers.
	args := []string{"container", "stop", "-t", "10"}
	for _, cntr := range containers {
		args = append(args, cntr.ID)
	}
	cmd := exec.Command(docker, args...)

	// Execute stop command.
	var out []byte
	if out, err = cmd.CombinedOutput(); nil != err {
		err = fmt.Errorf(
			"while running 'docker container stop': %s | %v", string(out), err,
		)
		return
	}

	// Prune wiped containers.
	err = Container.Prune()
	return
}
func (*containerNS) Prepare(js jumpServer.JumpServer, publicKeyBytes []byte) (cmd *exec.Cmd) {
	//to do remove this its just for testing purposes
	Container.Wipe()
	// Assign base arguments.
	p := fmt.Sprintf("%v", js.JumpSSHPort)
	sp := fmt.Sprintf("%v", js.SSHPort)
	vp := fmt.Sprintf("%v", js.VNCPort)
	bp := fmt.Sprintf("%v", js.BrowserPort)
	id_rsa_pub_path := "/tmp/" + p + "_id_rsa.pub"
	os.Remove(id_rsa_pub_path)
	ioutil.WriteFile(id_rsa_pub_path, publicKeyBytes, 0600)
	args := []string{
		"run", "-t", "-p", p + ":2222", "-v", id_rsa_pub_path + ":/root/.ssh/authorized_keys:ro",
	}
	if sp != "0" {
		args = append(args, "-p", sp+":"+sp)
	}
	if vp != "0" {
		args = append(args, "-p", vp+":"+vp)
	}
	if bp != "0" {
		args = append(args, "-p", bp+":"+bp)
	}

	// Assign Docker Image to execute.
	args = append(args, "rmd:latest")
	fmt.Println("starting docker with ", args)
	// Create the execution command.
	cmd = exec.Command(docker, args...)

	// Put the process in its own process group. This prevents SIGINT from
	// interrupting the running container so we can shut it down in a controlled
	// manner rather than have it shutdown, restart and be shut down, again.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	return

}
