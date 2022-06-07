package su

import "os"

// Privileges returns true if the user has super-user privilege.
func Privileges() bool {
	user := os.Getenv("USER")
	return user == "root"
}
