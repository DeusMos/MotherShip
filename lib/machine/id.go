package machine

import (
	"fmt"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/config"
)

func GetMachineID() (id string, err error) {
	hostName, err := getHostname()
	customer := config.Get.Reporting.Customer
	plant := config.Get.Reporting.Plant
	id = fmt.Sprintf("%v_%v_%v", customer, plant, hostName)

	return
}
