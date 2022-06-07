package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/jumpServer"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/machine"
)

func CheckPortal() (active bool, err error) {
	var myid string
	var js jumpServer.JumpServer
	endpoint := "http://woodfam.us:8888/api/getnextmessage"
	var client = http.Client{}
	myMachine, err := machine.Discover()

	if nil != err {
		fmt.Print(err)
		return
	}
	myid = fmt.Sprintf("%v_%v_%v", myMachine.Customer, myMachine.Plant, myMachine.Hostname)
	// session, err := ioutil.ReadFile("session")
	if err != nil {
		log.Fatal(err)
	}
	cred := func(r *http.Request) {
		// r.Header.Set("Fini-Session", string(session))
		r.Header.Set("Target-Machine", string(myid))
	}
	var req *http.Request

	if req, err = http.NewRequest("GET", endpoint, nil); nil != err {
		err = fmt.Errorf(
			"while requesting connection got: %v", err,
		)
		return
	}

	cred(req)

	b, _ := httputil.DumpRequest(req, false)
	fmt.Print(string(b))
	var resp *http.Response
	if resp, err = client.Do(req); nil != err {
		err = fmt.Errorf("error: while connecting got: %v", err)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&js)
	if nil != err {
		return
	}

	fmt.Printf("jumpserver: %+v ", js)
	if js.JumpSSHPort != 0 {
		active, err = js.Connect()
		if nil != err {
			return
		}
	}

	return

}

// opt are CLI-passed options. Each option should have a variable and a
// description of that variable for help documentation.
var opt struct {
	config     string
	configDesc string
}

// init sets up the CLI parser.
func Init() {
	opt.configDesc = "Path to optional configuration file."

	flag.StringVar(&opt.config, "config", "", opt.configDesc)
	flag.StringVar(&opt.config, "c", "", opt.configDesc)
}
func main() {

	Init()
	log.Println("Starting...")
	flag.Parse()
	MainLoop()
}

func MainLoop() {
	var active bool
	var err error
	active = false
	for {
		if !active {
			active, err = CheckPortal()
			if err != nil {
				fmt.Println(err)
			}
		}
		time.Sleep(15 * time.Second)
	}
}
