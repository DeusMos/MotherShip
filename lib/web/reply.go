package web

import (
	"encoding/json"
	"io"
	"net/http"
)

// Reply to the client with a provided message.
//   w: Writer used to send a reply back to the client.
//   status: HTTP status code for the response (ex, '200 OK').
//   obj: Response to be inserted into the body. May be an error, string,
//        io.Reader or JSON-serializable object.
//   cookies: Array of HTTP cookies to send to the client.
type Reply struct {
	w       http.ResponseWriter
	status  int
	obj     interface{}
	cookies []*http.Cookie
}

// Status code to be returned for the request.
func (reply *Reply) Status(status int) *Reply {
	reply.status = status
	return reply
}

// With the supplied message, reply to the client.
func (reply *Reply) With(obj interface{}) *Reply {
	reply.obj = obj
	return reply
}

// Cookie sets a cookie in the response.
func (reply *Reply) Cookie(cookie *http.Cookie) *Reply {
	reply.cookies = append(reply.cookies, cookie)
	return reply
}

// Forward a response.
func (reply *Reply) Forward(resp *http.Response) (err error) {
	// Copy all HTTP headers from the response.
	reply.w.WriteHeader(resp.StatusCode)
	for k, vs := range resp.Header {
		for _, v := range vs {
			reply.w.Header().Add(k, v)
		}
	}

	// Copy the contents of the response.
	_, err = io.Copy(reply.w, resp.Body)
	return
}

// Do performs the actual reply.
func (reply *Reply) Do() (err error) {
	// Add cookies
	for _, cookie := range reply.cookies {
		http.SetCookie(reply.w, cookie)
	}

	reply.w.WriteHeader(reply.status)

	if nil != reply.obj {

		// If the response object is an error, package it into the JSON layout
		// expected by the client.
		if err, ok := reply.obj.(error); ok {
			reply.obj = &errorMsg{
				Error: err.Error(),
			}
		}

		// If the response object is a string, send the string as plain text.
		if str, ok := reply.obj.(string); ok {
			reply.w.Header().Set("Content-type", "text/plain; charset=utf-8")
			_, err = io.WriteString(reply.w, str)
			return
		}

		// If a reader was passed in, directly copy the reader.
		if reader, ok := reply.obj.(io.Reader); ok {
			_, err = io.Copy(reply.w, reader)
			return
		}

		// The passed object will be serialized into JSON before being sent to
		// the client.
		reply.w.Header().Set("Content-type", "application/json")
		encoder := json.NewEncoder(reply.w)
		return encoder.Encode(reply.obj)
	}

	return
}

// errorMsg is used to pack any errors into the format expected by the client.
type errorMsg struct {
	Error string `json:"error"`
}
