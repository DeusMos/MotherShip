package machine

import "fmt"

type Machine struct {
	Hostname string
	Customer string
	Plant    string
	Alive    bool
	ID       string
}

func (m *Machine) ToString() (returnString string) {
	returnString = fmt.Sprintf("%s, %s, %s, %t, %s", m.Hostname, m.Customer, m.Plant, m.Alive, m.ID)
	return
}
