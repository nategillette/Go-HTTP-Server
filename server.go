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
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	loggly "github.com/jamespearly/loggly"
)

type status struct {
	HTTP int
	Time time.Time
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/status", Index).
		Methods("GET").
		Schemes("http", "https")
	router.HandleFunc("/*", Fail).
		Methods("GET").
		Schemes("http", "https")
	router.HandleFunc("/*", BadMethod).
		Methods("HEAD", "POST", "PUT", "PATCH",
			"DELETE", "CONNECT", "OPTIONS", "TRACE").
		Schemes("http", "https")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func BadMethod(w http.ResponseWriter, r *http.Request) {
	var tag string
	tag = "BadMethod"

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	client := loggly.New(tag)
	Status := status{

		HTTP: http.StatusMethodNotAllowed,
		Time: time.Now(),
	}

	msg, _ := json.Marshal(Status)

	err := client.Send("info", "Method:"+r.Method+
		",IP:"+ip+",Path:"+r.RequestURI+string(msg))

	fmt.Println("err:", err)

	w.Write(msg)
}

func Fail(w http.ResponseWriter, r *http.Request) {
	var tag string
	tag = "BadPath"

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	client := loggly.New(tag)
	Status := status{

		HTTP: http.StatusNotFound,
		Time: time.Now(),
	}

	msg, _ := json.Marshal(Status)

	err := client.Send("info", "Method:"+r.Method+
		",IP:"+ip+",Path:"+r.RequestURI+string(msg))

	fmt.Println("err:", err)

	w.Write(msg)
}

func Index(w http.ResponseWriter, r *http.Request) {
	var tag string
	tag = "GoodQuery"

	client := loggly.New(tag)

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	Status := status{

		HTTP: http.StatusOK,
		Time: time.Now(),
	}

	msg, _ := json.Marshal(Status)

	err := client.Send("info", "Method:"+r.Method+
		",IP:"+ip+",Path:"+r.RequestURI+string(msg))

	fmt.Println("err:", err)

	w.Write(msg)
}
