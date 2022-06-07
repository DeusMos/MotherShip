package message

import (
	"sync"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/jumpServer"
)

var once sync.Once

type Que map[string]Message

var instance Que

func GetQue() Que {
	once.Do(func() { instance = make(Que) })
	return instance
}

func Add(m Message, machineID string) {
	once.Do(func() { instance = make(Que) })
	instance[machineID] = m
}
func Get(machineID string) (m Message) {
	once.Do(func() { instance = make(Que) })
	var ok bool
	m, ok = instance[machineID]
	if !ok {
		js := jumpServer.NewJumpServer("", "", "", 0, 0, 0, 0)
		m = NewMessage(js, false, machineID)
	}
	return
}
