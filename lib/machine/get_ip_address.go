package machine

import (
	"net"
	"strings"
)

func getIPAddresses() (addrs []string, err error) {
	// Lookup all of the network interfaces.
	var interfaces []net.Interface
	if interfaces, err = net.Interfaces(); nil != err {
		return
	}

	// Search for the 'eth1' network interface.
	for _, inter := range interfaces {
		if "eth1" != inter.Name {
			continue
		}

		// Lookup all IP addresses on the 'eth1' network interface.
		var addresses []net.Addr
		if addresses, err = inter.Addrs(); nil != err {
			return
		}

		// Filter out IPv6 addresses.
		for _, addr := range addresses {
			a := addr.String()
			if strings.Contains(a, ":") {
				continue
			}

			// Add the desired IPv4 addressess from 'eth1' network interface.
			addrs = append(addrs, a)
		}
	}
	return
}
