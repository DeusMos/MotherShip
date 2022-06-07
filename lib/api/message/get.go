package message

import (
	"errors"
	"net/http"
	"net/http/httputil"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/message"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/web"
)

func Get(req *web.Req) {
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
	msg := message.Get(targetMachine)
	req.Reply.Status(http.StatusOK)
	req.Reply.With(msg.Js)
	req.Reply.Do()
}
