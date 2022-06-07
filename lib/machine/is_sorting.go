package machine

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func isSorting() (sorting bool, err error) {
	var contents []byte
	if contents, err = ioutil.ReadFile("/home/user/.g6/sortstate"); nil != err {
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

	// Verify the key is the expected 'value' key.
	key := bytes.TrimSpace(parts[0])
	if 0 != bytes.Compare([]byte("value"), key) {
		err = fmt.Errorf("expected a key of 'value' but got %s", string(key))
		return
	}

	value := string(bytes.TrimSpace(parts[1]))
	sorting = value == "1"
	return
}
