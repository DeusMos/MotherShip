package ssh

import (
	"errors"
	"net/http"
	"net/http/httputil"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/docker"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/id_rsa"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/jumpServer"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/message"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/web"
)

// GetReq to update system clock with the seconds since epoch.

func Get(req *web.Req) {
	// Parse the JSON request.
	b, _ := httputil.DumpRequest(req.R, false)
	targetMachine := req.R.Header.Get("Target-Machine")
	print("req.R = ", string(b), "targetMachine = ", targetMachine)
	var err error

	// Verify value was supplied for targetMachine.
	if targetMachine == "" {
		err = errors.New("'TargetMachine' not set in request")
		req.Reply.Status(http.StatusBadRequest).With(err).Do()
		return
	}
	publicKeyBytes, privateKeyBytes, err := id_rsa.GetNew()
	if err != nil {
		return
	}
	//to do look up lowest open port index
	js := jumpServer.NewJumpServer("woodfam.us", string(privateKeyBytes), string(publicKeyBytes), 2201, 10001, 0, 0)

	req.Reply.Status(http.StatusCreated)
	req.Reply.With(js)
	message.Add(message.NewMessage(js, true, targetMachine), targetMachine)

	c := docker.Container.Prepare(js, publicKeyBytes)
	// c.Stdout = newLogger(log.InfOut)//what do with logs?
	// c.Stderr = newLogger(log.InfErr)

	// Start the jumpserver.
	if err = c.Start(); nil != err {
		c = nil
		return
	}
	// Success!
	req.Reply.Do()

}
