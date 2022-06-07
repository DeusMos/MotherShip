package fini

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/jumpServer"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/session"
)

func GetJumpServer(endpoint string, target string) (js jumpServer.JumpServer, err error) {
	var client = http.Client{}
	s, err := session.Load()
	if err != nil {
		fmt.Print("unable to find session did you log in?")
		return
	}
	cred := func(r *http.Request) {
		r.Header.Set("Fini-Session", s.Session)
		r.Header.Set("Target-Machine", string(target))
	}
	var req *http.Request

	if req, err = http.NewRequest("GET", endpoint, nil); nil != err {
		err = fmt.Errorf(
			"while requesting connection got: %v", err,
		)
		return
	}
	cred(req)
	b, err := httputil.DumpRequest(req, false)
	fmt.Println(string(b))
	var resp *http.Response
	if resp, err = client.Do(req); nil != err {
		err = fmt.Errorf("while connecting got: %v", err)
		return
	}
	defer resp.Body.Close()
	b, err = httputil.DumpResponse(resp, true)
	if nil != err {
		fmt.Print(err.Error())
	}
	fmt.Println(string(b))
	err = json.NewDecoder(resp.Body).Decode(&js)
	if nil != err {
		fmt.Print(err.Error())
	}
	return
}
