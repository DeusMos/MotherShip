package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SortGroup ...
type SortGroup struct {
	Keyware    string   `json:"keyware"`
	Settings   string   `json:"settings"`
	SensorType string   `json:"sensorType"`
	Sensors    []uint16 `json:"sensors"`
}

// Inferencer ...
type Inferencer struct {
	Running bool   `json:"running"`
	Version string `json:"version"`
	Online  bool   `json:"online"`
}

// Message ...
type Message struct {
	Customer    string       `json:"customer"`
	ID          string       `json:"id"`
	Plant       string       `json:"plant"`
	Hostname    string       `json:"hostname"`
	IPAddresses []string     `json:"ipAddresses"`
	Product     string       `json:"product"`
	SortGroups  []*SortGroup `json:"sortGroups"`
	Sorting     bool         `json:"sorting"`
	Version     string       `json:"version"`
}

// Marshal ...
func (m *Message) Marshal() (d []byte, err error) {
	return json.MarshalIndent(m, "", "  ")
	// return json.Marshal(m)
}

// Unmarshal ...
func (m *Message) Unmarshal(d []byte) (err error) {
	return json.Unmarshal(d, m)
}

// SaveToFile ...
func (m *Message) SaveToFile(filename string) (err error) {
	// Get the bytes.
	var contents []byte
	if contents, err = m.Marshal(); nil != err {
		return
	}

	// TODO: REMOVE CODE
	fmt.Println(string(contents))

	// TODO: Store the bytes.
	// err = ioutil.WriteFile(filename, contents, 0644)
	return
}

// UploadToServer ...
func (m *Message) UploadToServer(endpoint string) (err error) {
	var contents []byte
	if contents, err = m.Marshal(); nil != err {
		return
	}
	client := &http.Client{}
	var resp *http.Response
	buffer := bytes.NewBuffer(contents)
	var req *http.Request
	if req, err = http.NewRequest("POST", endpoint, buffer); nil != err {
		err = fmt.Errorf("while Posting machine data to server got: %v", err)
		return
	}
	req.Header.Set("Fini-Machine-Type", "g6")
	req.Header.Add("Content-Type", "application/json")
	session, err := ioutil.ReadFile("./session")
	if nil != err {
		err = fmt.Errorf("while loading session got: %v", err)
		return
	}
	cookie := &http.Cookie{Name: "FiniSession", Value: string(session), HttpOnly: false}
	req.AddCookie(cookie)

	if resp, err = client.Do(req); nil != err {
		return
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		var msg []byte
		if msg, err = ioutil.ReadAll(resp.Body); nil != err {
			return
		}
		err = fmt.Errorf(
			"http request returned %d: %s: %s",
			resp.StatusCode, http.StatusText(resp.StatusCode), string(msg),
		)
		return
	}

	return
}
