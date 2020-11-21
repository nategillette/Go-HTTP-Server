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
