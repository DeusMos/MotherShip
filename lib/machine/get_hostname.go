package machine

import (
	"bytes"
	"io/ioutil"
)

func getHostname() (name string, err error) {
	var contents []byte
	if contents, err = ioutil.ReadFile("/etc/hostname"); nil != err {
		return
	}

	name = string(bytes.TrimSpace(contents))
	return
}
