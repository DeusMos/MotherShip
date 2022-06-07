package machine

import (
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/config"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/data"
)

// Discover ...
func Discover() (info *data.Message, err error) {
	info = &data.Message{
		ID:       config.Get.Reporting.ID,
		Customer: config.Get.Reporting.Customer,
		Plant:    config.Get.Reporting.Plant,
	}

	// Get Hostname.
	if info.Hostname, err = getHostname(); nil != err {
		return
	}

	// // Get IP Addresses.
	// if info.IPAddresses, err = getIPAddresses(); nil != err {
	// 	return
	// }

	// // Get Product name.
	// if info.Product, err = getProduct(); nil != err {
	// 	return
	// }

	// // Get sort groups.
	// if info.SortGroups, err = getSortGroups(); nil != err {
	// 	return
	// }

	// // Get sorting status.
	// if info.Sorting, err = isSorting(); nil != err {
	// 	return
	// }
	return
}
