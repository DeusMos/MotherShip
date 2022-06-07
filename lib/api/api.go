package api

import (
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/api/message"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/api/ssh"
)

var (
	SSH     ssh.Struct
	Message message.Struct
)

func init() {
	SSH = ssh.API
	Message = message.API
}
