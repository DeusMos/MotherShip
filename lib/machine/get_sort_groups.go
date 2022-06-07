package machine

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/data"
)

func getSortGroups() (groups []*data.SortGroup, err error) {
	// Read contents for sortmap file.
	var contents []byte
	if contents, err = ioutil.ReadFile("/home/user/.g6/sortmap"); nil != err {
		return
	}

	// Split the contents into each line entry.
	lines := bytes.Split(bytes.TrimSpace(contents), []byte("\n"))

	// Parse each line entry.
	for _, line := range lines {
		var sg *data.SortGroup
		if sg, err = parseSortGroupLine(line); nil != err {
			return
		}

		groups = append(groups, sg)
	}
	return
}

// parseSortGroupLine takes an entire line entry and parses it into a SortGroup.
func parseSortGroupLine(line []byte) (sg *data.SortGroup, err error) {
	// Split the sort group line into its key:value pairs.
	// '0,1,2	frozenPotatoStripSLC360EFT Camera 1/4 SkinOn' ->
	// ['0,1,2', 'frozenPotatoStripSLC360EFT Camera 1/4 SkinOn']
	parts := bytes.SplitN(bytes.TrimSpace(line), []byte("\t"), 2)
	if 2 != len(parts) {
		err = fmt.Errorf(
			"while getting sort group expected key and value but got %d parts",
			len(parts),
		)
		return
	}

	// Assign key ('0,1,2').
	key := bytes.TrimSpace(parts[0])
	// Assign value ('frozenPotatoStripSLC360EFT Camera 1/4 SkinOn').
	value := bytes.TrimSpace(parts[1])

	// Parse the sort group information.
	tmp := &data.SortGroup{}
	if tmp.Sensors, tmp.SensorType, err = parseSensors(key); nil != err {
		return
	}
	if tmp.Keyware, tmp.Settings, err = parseKeyware(value); nil != err {
		return
	}

	sg = tmp
	return
}

// parseSensors for string-encoded byte slice to slice of sensor indexes and
// their sensor types.
func parseSensors(value []byte) (sensors []uint16, sType string, err error) {
	// Split out the numbers based on commas '0,1,2' -> ['0', '1', '2'].
	values := bytes.Split(value, []byte(","))

	for _, vByte := range values {
		// Parse index from byte slice -> string -> int
		// ([]byte("1") -> "1" -> 1).
		vStr := string(vByte)
		var n int
		if n, err = strconv.Atoi(vStr); nil != err {
			err = fmt.Errorf(
				"while parsing sensors expected a number but got %s", vStr,
			)
			return
		}

		// Verify the sensor index is sane.
		if n < 0 || n > 16 {
			err = fmt.Errorf(
				"expected a small positive integer for sensor but got %d", n,
			)
			return
		}

		// Add the sensor index to the return slice.
		sensors = append(sensors, uint16(n))
	}

	// If there are no sensor indexes, return an error.
	if 0 == len(sensors) {
		err = errors.New("no sensors found in sensor group")
		return
	}

	// Get the sensor type for each sensor index and verify that all sensors are
	// the same sensor type.
	for _, s := range sensors {

		// Get the sensor type and verify it is valid.
		tmpType := getSensorType(s)
		if "" == tmpType {
			err = fmt.Errorf(
				"couldn't determine sensor type for sensor index %d", s,
			)
			return
		}

		if "" == sType {
			// If the sensor group type has not been set, update it.
			sType = tmpType
		} else {
			// If the sensor group type has been set, verify this sensor type
			// matches that of the group.
			if tmpType != sType {
				err = fmt.Errorf(
					"sensor type did not match group (%s vs %s for %d)",
					sType, tmpType, s,
				)
				return
			}
		}
	}

	return
}

// parseKeyware value into keyware name and saved settings name.
func parseKeyware(value []byte) (keyware, settings string, err error) {
	// Split 'frozenPotatoStripSLC360EFT Camera 1/4 SkinOn' into
	// ['frozenPotatoStripSLC360EFT', 'Camera 1/4 SkinOn'] and validate. It is
	// possible for there to only be one part (no saved settings).
	parts := bytes.SplitN(value, []byte(" "), 2)
	if 2 < len(parts) || 0 == len(parts) {
		err = fmt.Errorf(
			"while parsing keyware expected 1 or 2 parts but got %d",
			len(parts),
		)
		return
	}

	// Assign the keyware ('frozenPotatoStripSLC360EFT').
	keyware = string(bytes.TrimSpace(parts[0]))

	// If the saved settings exist, assign settings ('Camera 1/4 SkinOn').
	if 2 == len(parts) {
		settings = string(bytes.TrimSpace(parts[1]))
	}

	return
}

// getSensorType for the passed sensor index.
func getSensorType(index uint16) (sType string) {
	if exists(fmt.Sprintf("/home/user/.g6/is%d", index)) {
		// Is the sensor a camera?
		sType = "camera"
		return

	} else if exists(fmt.Sprintf("/home/user/.g6/ls%d", index)) {
		// Is the sensor a laser?
		sType = "laser"
		return
	}
	return
}

// exists checks to see if a file or directory exists.
func exists(filename string) bool {
	if _, err := os.Stat(filename); nil != err {
		return false
	}
	return true
}
