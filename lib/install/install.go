package install

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gitlab.rd.keyww.com/sw/spike/rmd/cmd/clitools"
	certs "gitlab.rd.keyww.com/sw/spike/rmd/lib/certs/content"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/machine"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/su"
)

var initTabLine = "ktin:2:respawn:/usr/bin/rmd -c /etc/rmd/config.yml >/var/log/rmd.log 2>&1"

// Install the inferencer.
func Install() (err error) {

	// Verify root permissions.
	if !su.Privileges() {
		err = errors.New("install requires root permissions")
		return
	}

	// Verify we are located in the proper location.
	expected := "/usr/bin/rmd"
	if expected != os.Args[0] {
		err = fmt.Errorf(
			"rmd is located in '%s' but should be in '%s'",
			os.Args[0], expected,
		)
		return
	}

	_ = os.Mkdir("/etc/ssl", os.ModePerm)

	_ = os.Mkdir("/etc/ssl/certs", os.ModePerm)
	// If ca-certificates.crt does not exist, add it.
	if !exists(certs.Path) {
		if err = ioutil.WriteFile(
			certs.Path, []byte(certs.Contents), 0644,
		); nil != err {
			return
		}
	}
	fmt.Println("getting info")
	info, err := machine.Discover()
	fmt.Println(info)
	if nil != err {
		return
	}
	fmt.Println("registering machine")
	clitools.RegisterMachine(info.Customer + "_" + info.Plant + "_" + info.Hostname)
	// If inference is not in /etc/inittab, add it or update it.
	if err = updateInitTab(); nil != err {
		return
	}

	return
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return nil == err
}

func updateInitTab() (err error) {
	var contents []byte
	if contents, err = ioutil.ReadFile("/etc/inittab"); nil != err {
		return
	}
	exists := bytes.Contains(contents, []byte("/usr/bin/rmd"))
	current := bytes.Contains(contents, []byte("/var/log/rmd.log"))
	if exists && current {
		return
	}
	if exists && !current {
		if contents, err = replaceInitTabLine(contents); nil != err {
			return
		}
	} else {
		parts := [][]byte{contents, []byte(initTabLine)}
		contents = bytes.Join(parts, []byte("\n"))
	}

	if err = ioutil.WriteFile("/etc/inittab", contents, 0644); nil != err {
		return
	}

	cmd := exec.Command("telinit", "q")
	var out []byte
	if out, err = cmd.Output(); nil != err {
		err = fmt.Errorf("got '%v' with: %s", err, string(out))
		return
	}

	return
}

func replaceInitTabLine(in []byte) (out []byte, err error) {
	parts := bytes.Split(in, []byte("\n"))
	for i, part := range parts {
		if bytes.Contains(part, []byte("/usr/bin/rmd")) {
			parts[i] = []byte(initTabLine)
			out = bytes.Join(parts, []byte("\n"))
			return
		}
	}
	err = errors.New("line to replace not found")
	return
}
