package ssh

import "gitlab.rd.keyww.com/sw/spike/rmd/lib/web"

// Struct defining API tree.
type Struct struct {
	Get web.Handle
}

var (
	// API tree.
	API Struct
)

func init() {
	API.Get = Get
}
