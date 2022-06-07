package web

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handle is any function that accepts a Req.
type Handle func(*Req)

// Req (uest) encapsulates the request information and creates a reply upon
// request.
type Req struct {
	W     http.ResponseWriter
	R     *http.Request
	PS    httprouter.Params
	Reply *Reply
}

// NewReq creates a new HTTP connection request.
func NewReq(
	w http.ResponseWriter, r *http.Request, ps httprouter.Params,
) (req *Req) {
	req = &Req{
		W:  w,
		R:  r,
		PS: ps,
		Reply: &Reply{
			w:      w,
			status: http.StatusOK,
		},
	}
	return
}

// Decode encoded request based on content type.
func (req *Req) Decode(v interface{}) error {
	// Always assumes JSON, for now.
	decoder := json.NewDecoder(req.R.Body)
	return decoder.Decode(v)
}
