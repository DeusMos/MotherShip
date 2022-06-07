package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/api"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/config"
	"gitlab.rd.keyww.com/sw/spike/rmd/lib/web"
)

var (
	router *httprouter.Router
	server *http.Server
)

// Start the HTTP server.
func Start() (err error) {
	// If the server is running, stop it.
	if nil != server {
		Stop()
	}

	// Setup a new router and assign the API endpoints.
	router = httprouter.New()
	if err = setupAPI(); nil != err {
		return
	}

	// Start the HTTP server.
	err = startServer()
	return
}

// Stop the HTTP server.
func Stop() {
	// Give the server up to 1 minute to shutdown.
	func() {
		ctx, cancel := context.WithTimeout(
			context.Background(), time.Minute,
		)
		defer cancel()

		server.Shutdown(ctx)
	}()

	// Mark the server as stopped.
	server = nil
}

func setupAPI() (err error) {
	// Make the following code a bit more concise.
	get := router.GET
	//post := router.POST
	// put := router.PUT
	wrap := web.Wrap

	// Setup the API endpoints.
	get("/api/ssh", wrap(api.SSH.Get))
	get("/api/getnextmessage", wrap(api.Message.Get))
	// put("/api/firmware", wrap(api.Firmware.Put))
	// post("/api/load", wrap(api.Load.Post))
	// get("/api/logs", wrap(api.Logs.Get))
	// put("/api/model", wrap(api.Model.Put))
	// get("/api/status", wrap(api.Status.Get))

	return
}

// startServer starts the HTTP server.
func startServer() (err error) {
	// Create the server instance with proper binding (likely ':80').
	bindAddr := fmt.Sprintf(
		"%s:%d", config.Get.Server.Host, config.Get.Server.Port,
	)
	server = &http.Server{
		Addr:    bindAddr,
		Handler: router,
	}

	// Check if an error is thrown during starting of the HTTP server. If no
	// error is immediately thrown, continue operations.
	errC := make(chan error, 1)
	go func() {
		errC <- server.ListenAndServe()
	}()

	// Wait up to 2 seconds for the server to fail to start.
	timer := time.NewTimer(time.Second * 2)
	select {
	case err = <-errC:
	case <-timer.C:
	}

	return
}
