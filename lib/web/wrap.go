package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/julienschmidt/httprouter"
)

// Wrap a handler to ensure the request body is closed.
func Wrap(handler Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		b, err := httputil.DumpRequest(r, false)
		if err != nil {
			return
		}
		fmt.Print(string(b))
		fmt.Print("\nw=",w,"\nr=", r,"\nps = ", ps)

		handler(NewReq(w, r, ps))
		r.Body.Close()
	}
}
