package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	keys  = "keys"
	login = "login"
	pass  = "pass"
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type Answers struct {
	Not_send []string `json:"not_send"`
	Send     int      `json:"send"`
}

type Phone struct {
	CRM []struct {
		ToPhone string `json:"tophone"`
		Id      string `json:"id"`
	}
	Body      string `json:"body"`
	Fromname  string `json:"fromname"`
	Key_phone string `json:"key_phone"`
}

func sms(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Request-Headers", "*")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	body, err := ioutil.ReadAll(req.Body)
	checkError(err)
	var phone Phone
	var answers Answers
	json.Unmarshal([]byte(body), &phone)
	if phone.Key_phone == keys {
		for _, crm := range phone.CRM {
			res_answer := sendphone(crm.ToPhone, phone.Body, phone.Fromname)
			if res_answer == "0" {
				answers.Not_send = append(answers.Not_send, crm.Id)
			} else {
				answers.Send = answers.Send + 1
			}
		}
	}
	ret, _ := json.Marshal(answers)
	w.WriteHeader(200)
	w.Write([]byte(ret))
}

func main() {
	http.HandleFunc("/", sms)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendphone(tophone string, body_s string, fromname string) (answer string) {
	xmlBody := `<?xml version="1.0" encoding="utf-8"?>
	              <package login="` + login + `" password="` + pass + `">
	               <message>
		            <default sender="` + fromname + `"/>
		             <msg recipient="` + tophone + `">` + body_s + `</msg>
	               </message>
	              </package>`

	resp, _ := http.Post("https://xmlapi.devinotele.com/Send.ashx", "text/xml", strings.NewReader(xmlBody))
	body, _ := ioutil.ReadAll(resp.Body)

	if strings.Contains((string(body)), "106</msg>") {
		answer = "1"
	} else {
		answer = "0"
	}
	resp.Body.Close()
	return answer
}
