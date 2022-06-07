package service

import (
	"fmt"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/service/http"
)

// Start all services.
func Start() (err error) {
	msg := "while starting '%s' service: %v"
	if err = http.Start(); nil != err {
		err = fmt.Errorf(msg, "HTTP API", err)
		return
	}

	return
}
func Stop() {
	http.Stop()
}
