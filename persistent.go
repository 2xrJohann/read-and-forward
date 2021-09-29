package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"bufio"
	"strings"

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

//wait for input on stdin, if valid json then make the request, otherwise discard
func waitForInput(){
	validate := validator.New()
	ch := make(chan string)
	done := make(chan interface{})
	reader := bufio.NewReader(os.Stdin)

	go func(ch chan string, reader *bufio.Reader, done chan interface{}){
		for{
			text, _ := reader.ReadString('\n')
			text = strings.TrimSuffix(text, "\n")

			/*windows has ascii char 13 on the end of strings
			fir carriage returns when I use vscode, if its there
			then get rid of that too since it messes up my done
			channel*/
			ascii := 13
			windowsTrim := string(ascii)
			text = strings.TrimSuffix(text, windowsTrim)

			if strings.Compare(text, "done") == 0{ 
				notification := new(interface{})
				done <- notification
			}else{
				ch <- text
			}
		}
	}(ch, reader, done)

	for{
		select{
		case data:= <- ch:
			var input inputStructure
			if err := json.Unmarshal([]byte(data), &input); err != nil {
				fmt.Println(err)
			}else{
				err := validate.Struct(input)
				if err != nil {
					fmt.Println(err)
				}else{
					go sendRequest(&input)
				}
			}
		case <- done:
			return
		}
	}
}

func main() {

	//Although not required, listen on port 8080 for requests on index url for proof of concept
	reqHandler := func(w http.ResponseWriter, req *http.Request) {
		//handler just for proof of concept, do nothing
	}

	http.HandleFunc("/", reqHandler)
	go http.ListenAndServe(":8080", nil)

	waitForInput()
}
