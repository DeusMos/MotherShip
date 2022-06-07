package message

import (
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/jumpServer"
)

type Message struct {
	Js        jumpServer.JumpServer
	Active    bool
	MachineID string
}

func NewMessage(js jumpServer.JumpServer, active bool, machineID string) Message {
	return Message{Js: js, Active: active, MachineID: machineID}
}
