package log

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gitlab.rd.keyww.com/sw/spike/rmd/lib/config"
)

// MgrOut is used to log a Mgr message to stdout.
func MgrOut(v ...interface{}) {
	entry := fmt.Sprint(v...)
	write(mgr, out, entry)
}

// MgrOutf is used to log a formatted Mgr message to stdout.
func MgrOutf(format string, v ...interface{}) {
	entry := fmt.Sprintf(format, v...)
	write(mgr, out, entry)
}

// MgrErr is used to log a Mgr message to stderr.
func MgrErr(v ...interface{}) {
	entry := fmt.Sprint(v...)
	write(mgr, errr, entry)
}

// MgrErrf is used to log a formatted Mgr message to stderr.
func MgrErrf(format string, v ...interface{}) {
	entry := fmt.Sprintf(format, v...)
	write(mgr, errr, entry)
}

// InfOut is used to log a Inference message to stdout.
func InfOut(v ...interface{}) {
	entry := fmt.Sprint(v...)
	write(infer, out, entry)
}

// InfOutf is used to log a formatted Inference message to stdout.
func InfOutf(format string, v ...interface{}) {
	entry := fmt.Sprintf(format, v...)
	write(infer, out, entry)
}

// InfErr is used to log a Inference message to stderr.
func InfErr(v ...interface{}) {
	entry := fmt.Sprint(v...)
	write(infer, errr, entry)
}

// InfErrf is used to log a formatted Inference message to stderr.
func InfErrf(format string, v ...interface{}) {
	entry := fmt.Sprintf(format, v...)
	write(infer, errr, entry)
}

// History returns a single string that contains all log entries.
func History() string {
	// If there is not history, inform the requestor.
	if !toHistory {
		return "LOGGING HISTORY IS NOT ENABLED!"
	}

	// Create and submit a history request and wait for it to be populated.
	// Return the populated value.
	h := &historyReq{
		wg: &sync.WaitGroup{},
	}
	h.wg.Add(1)
	historyC <- h
	h.wg.Wait()
	return h.history
}

// Start the logging service.
func Start() (err error) {
	// If console logging is enabled, set up console logging.
	if config.Get.Logging.EnableConsole {
		toConsole = true
	}

	// If history logging is enabled, set up history logging.
	if config.Get.Logging.EnableHistory {
		toHistory = true
		max = config.Get.Logging.HistoryLimit
		full = false
		ptr = 0
		history = make([]string, max)
	}

	// Start the logging service.
	wg.Add(1)
	go run()
	return
}

// Stop the logging service.
func Stop() {
	stopC <- struct{}{}
	wg.Wait()
}

// Constants used for consistent logging.
//   mgr: Logs originating from the Manager service.
//   infer: Logs originating from the Inference service.
//   out: Logs destined for stdout.
//   err: Logs destined for stderr.
const (
	mgr   = "MGR"
	infer = "INFER"
	out   = "OUT"
	errr  = "ERR"
)

// Variables for tracking how logging will be accomplished.
//   toConsole: True if logs should be printed to the console.
//   toHistory: True if logs should be tracked in history.
//   history: List of log entries.
//   full: True after the rotating log has filled the first pass through.
//   ptr: Pointer to the next entry to be (over)written.
//   max: Max number of history entries to track.
var (
	toConsole bool
	toHistory bool
	history   []string
	full      bool
	ptr       uint32
	max       uint32
)

// Variables for managing the service.
//   wg: Unblocks when logging service has shutdown.
//   stopC: Channel used to trigger shutdown of logging service.
//   msgC: Channel used to add new log entries.
//   historyC: Channel used to request logging history.
var (
	wg       = &sync.WaitGroup{}
	stopC    = make(chan struct{})
	msgC     = make(chan string, 16)
	historyC = make(chan *historyReq, 2)
)

// historyReq is used to send a request for the history into the logging routine
// and wait for the history value to be populated.
type historyReq struct {
	history string
	wg      *sync.WaitGroup
}

// handleHistory collects all logs and responds to the history request.
func handleHistory(req *historyReq) {
	// Inform the requestor when we are done collecting the history.
	defer req.wg.Done()

	// Determine the number of history entries that will be sent.
	var coll []string
	if full {
		coll = make([]string, max)
	} else {
		coll = make([]string, ptr)
	}

	// If history is full, get the entries after the point first, as these will
	// be the oldest entries.
	i := 0
	if full {
		for h := ptr; h < max; h++ {
			coll[i] = history[h]
			i++
		}
	}

	// Now gather the entries before the pointer.
	for h := uint32(0); h < ptr; h++ {
		coll[i] = history[h]
		i++
	}

	// Join the entries into a single string. This also allocates a new string
	// so any access to this string after it is created is race-safe.
	req.history = strings.Join(coll, "\n")
}

// handleMsg logging to desired outputs.
func handleMsg(msg string) {
	// Log to the console, if desired.
	if toConsole {
		fmt.Println(msg)
	}

	// Log to history, if desired. Increment to the next history entry position.
	// If we run off the end of the history, restart from the beginning.
	if toHistory {
		history[ptr] = msg
		ptr++
		if ptr >= max {
			ptr = 0
			full = true
		}
	}
}

// run the logging service.
func run() {
	defer wg.Done()

	for {
		select {
		case msg := <-msgC:
			// Handle a new log entry that needs routed.
			handleMsg(msg)

		case req := <-historyC:
			// Handle a history request.
			handleHistory(req)

		case <-stopC:
			// Shutdown the service.
			return
		}
	}
}

// write a message to be logged.
func write(service, stream, entry string) {
	now := time.Now()
	msg := fmt.Sprintf(
		"%04d/%02d/%02d %02d:%02d:%02d [%s][%s] %s",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		service, stream, entry,
	)
	msgC <- msg
}
