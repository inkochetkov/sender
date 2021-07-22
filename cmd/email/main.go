package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"gopkg.in/gomail.v2"
)

const (
	servername     = "servername"
	serverlogin    = "serverlogin"
	serverpassword = "serverpassword"
	keys           = "keys"
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type Email struct {
	CRM []struct {
		ToEmail string `json:"toEmail"`
		Id      string `json:"id"`
	}
	Subj      string `json:"subj"`
	Body      string `json:"body"`
	Key_email string `json:"key_email"`
	Fromname  string `json:"fromname"`
}

func emails(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Request-Headers", "*")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	body, err := ioutil.ReadAll(req.Body)
	checkError(err)
	var email Email
	var result []string
	json.Unmarshal([]byte(body), &email)
	if email.Key_email == keys {
		for _, crm := range email.CRM {
			status := sendmail(crm.ToEmail, email.Subj, email.Body, email.Fromname)
			if status == "не отправлено" {
				result = append(result, crm.Id)
			}
		}
	}
	ret, _ := json.Marshal(result)
	w.WriteHeader(200)
	w.Write([]byte(ret))
}

func main() {
	http.HandleFunc("/", emails)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func sendmail(toEmail string, subj string, body string, fromname string) (status string) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", serverlogin, fromname)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subj)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(servername, 25, serverlogin, serverpassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	status = "отправлено"
	if err := d.DialAndSend(m); err != nil {
		status = "не отправлено"
		exec.Command("exipick -i | xargs exim -Mrm").Output()
	}
	return status
}
