/**
* TODO:
* have the server listen to the /status endpoint
* respond to GET requests with a JSON body containing the system time and an HTTP status of 200.
* All requests to other endpoints should result in an HTTP status of 404
* all requests using anything other than the GET method should return an HTTP status of 405
* a record of all requests should be sent to Loggly including the method type, source IP address, request path, and the resulting HTTP status code
*
 */
package main

import (
	"fmt"
	"net/http"

	loggly "github.com/jamespearly/loggly"
)

func hello(w http.ResponseWriter, req *http.Request) {
	client := loggly.New("hello")
	fmt.Fprintf(w, "hello\n")
	client.Send("info", "Hello Request")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			client := loggly.New("header")
			fmt.Fprintf(w, "%v: %v\n", name, h)
			client.Send("info", "Headers Request")
		}
	}
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8090", nil)
}
