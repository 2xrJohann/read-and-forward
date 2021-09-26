package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
)

//Expected structure of JSON input as a struct
type inputStructure struct {
	TxId       string           `json:"TxId" validate:"required,numeric"`
	TxType     string           `json:"TxType" validate:"required"`
	Result     string           `json:"Result" validate:"required"`
	ResultCode string           `json:"Resultcode" validate:"required,numeric"`
	Messages   messageStructure `json:"Messages"`
}

//Expected format of a message object inside the JSON input
type messageStructure struct {
	Message string `json:"Response"`
	Error   string `json:"Error"`
}

/*
Open input file and unmarshal json into struct
Return a pointer to the struct
*/
func readInput(path string) *inputStructure {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	var input inputStructure
	json.Unmarshal(data, &input)

	return &input
}

//Helper function to format and encode basic auth params for requests
func basicAuth(username string, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//Format the struct as JSON and send it as POST data to localhost endpoint with basic auth header
func sendRequest(input *inputStructure) {
	reqBody, err := json.Marshal(*input)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth("username", "password"))

	//no response will be returned from this request so ignore response return of client.Do()
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var validate *validator.Validate
	//Although not required, listen on port 8080 for requests on index url for proof of concept
	reqHandler := func(w http.ResponseWriter, req *http.Request) {
		//handler just for proof of concept, do nothing
	}

	http.HandleFunc("/", reqHandler)
	go http.ListenAndServe(":8080", nil)

	//Validate JSON input as per struct validation tags
	validate = validator.New()
	input := readInput("input.json")
	err := validate.Struct(input)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", *input)
		sendRequest(input)
	}
}
