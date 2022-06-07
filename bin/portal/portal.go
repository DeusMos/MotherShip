package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/config"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/service"
	lg "gitlab.rd.keyww.com/sw/spike/rmd/lib/service/log"
)

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
	log.Println("Starting...")
	flag.Parse()
	Init()
	// Load configuration, if applicable.
	if opt.config != "" {
		if err := config.Load(opt.config); nil != err {
			log.Fatalf("While loading configuration: %v\n", err)
		}
	}

	// Listen for shutdown signal.
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	// Start all services.
	if err := service.Start(); nil != err {
		log.Fatalf("While starting services: %v\n", err)
	}
	lg.MgrOut("Running...")

	// Wait for shutdown request.
	<-stop
	lg.MgrOut("Stopping...")
	service.Stop()
	log.Println("Stopped.")
}
