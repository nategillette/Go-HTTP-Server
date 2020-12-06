package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	loggly "github.com/jamespearly/loggly"
)

type status struct {
	HTTP int
	Time time.Time
}

type Item struct {
	ID          string `json:"ID"`
	Title       string
	Center      string
	Description string
	URL         string
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ngillet2/all", Contents).
		Methods("GET").
		Schemes("http", "https")
	router.HandleFunc("/ngillet2/status", Status).
		Methods("GET").
		Schemes("http", "https")
	router.HandleFunc("/ngillet2/search", Query).
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

	fmt.Println(msg)

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
	fmt.Println(msg)

	w.Write(msg)
}

func Status(w http.ResponseWriter, r *http.Request) {

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	var tag string
	tag = "Status"

	client := loggly.New(tag)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := dynamodb.New(sess)

	var tableName = "ngillet2-NASAPhotos"

	describeTables := &dynamodb.DescribeTableInput{
		TableName: &tableName,
	}

	resp, err := svc.DescribeTable(describeTables)

	if err != nil {
		fmt.Println("err:", err)
	}

	count := int(*resp.Table.ItemCount)

	//fmt.Println(count)

	tbl := *resp.Table.TableName

	//fmt.Println(tbl)

	data := `[{"TableName":` + tbl + `, "ItemCount":` + strconv.Itoa(count) + `}]`

	msg := []byte(data)

	err = client.Send("info", "Method:"+r.Method+
		",IP:"+ip+",Path:"+r.RequestURI+string(msg))

	fmt.Println("err:", err)

	w.Write(msg)

}

func Contents(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	var tag string
	tag = "All"

	client := loggly.New(tag)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := dynamodb.New(sess)

	var tableName = "ngillet2-NASAPhotos"

	describeTables := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	resp, _ := svc.Scan(describeTables)

	message := []Item{}

	dynamodbattribute.UnmarshalListOfMaps(resp.Items, &message)

	msg, _ := json.Marshal(message)

	err = client.Send("info", "Method:"+r.Method+
		",IP:"+ip+",Path:"+r.RequestURI+string(msg))

	fmt.Println("err:", err)

	//fmt.Println(message)

	w.Write(msg)
}
func Query(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	var tag string
	tag = "Query"

	client := loggly.New(tag)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := dynamodb.New(sess)

	var tableName = "ngillet2-NASAPhotos"

	query := r.URL.Query()

	ID, present := query["ID"]
	if !present || len(ID) == 0 {
		fmt.Println("No IDs Present")
	}
	if len(ID) > 1 {
		fmt.Println("Too Many Query Params")
	}

	param := ID[0]

	found, err := regexp.MatchString("^.+$", param)
	fmt.Printf("found=%v, err=%v", found, err)

	queryInput := &dynamodb.QueryInput{
		TableName: &tableName,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":ID": {
				S: aws.String(param),
			},
		},
		KeyConditionExpression: aws.String("ID = :ID"),
	}

	resp, err := svc.Query(queryInput)

	if err != nil {
		fmt.Println("err:", err)
	}
	message := []Item{}

	dynamodbattribute.UnmarshalListOfMaps(resp.Items, &message)

	msg, _ := json.Marshal(message)

	err = client.Send("info", "Method:"+r.Method+
		",IP:"+ip+",Path:" /*+r.RequestURI+string(msg)*/)

	if err != nil {
		fmt.Println("err:", err)
	}
	w.Write(msg)
	//fmt.Println(msg)

}
