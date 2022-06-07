package machine

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func getProduct() (name string, err error) {
	// Get the contents of the file.
	var contents []byte
	if contents, err = ioutil.ReadFile("/home/user/.g6/product"); nil != err {
		return
	}

	// Divide the contents into a single key:value pair (2 parts).
	parts := bytes.SplitN(contents, []byte("\t"), 2)
	if 2 != len(parts) {
		err = fmt.Errorf(
			"expected to have key and value but got %d parts",
			len(parts),
		)
		return
	}

	// Verify the key is the expected 'name' key.
	key := bytes.TrimSpace(parts[0])
	if 0 != bytes.Compare([]byte("name"), key) {
		err = fmt.Errorf("expected a key of 'name' but got %s", string(key))
		return
	}

	// Clean and return the value.
	name = string(bytes.TrimSpace(parts[1]))
	return
}
