package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"net/url"
	"bytes"

)

import (
   "google.golang.org/appengine"
	 "google.golang.org/appengine/urlfetch"

)
// [START import]
import (
	"twiml"
	"twirest"
)

// [END import]

func main() {
	http.HandleFunc("/call/receive", receiveCallHandler)
	http.HandleFunc("/sms/send", sendSMSHandler)
	http.HandleFunc("/sms/receive", receiveSMSHandler)

	appengine.Main()
}

var (
	twilioClient = twirest.NewClient("ACCOUNT-SID", "AUTH-TOKEN")
	twilioNumber = "TWILIO-NUMBER"
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func receiveCallHandler(w http.ResponseWriter, r *http.Request) {
	resp := twiml.NewResponse()
	resp.Action(twiml.Say{Text: "Wolfram Alpha is the world's most advanced computational knowledge engine, but I can't receive voice queries yet. Send me a text message instead."})
	resp.Send(w)
}

func sendSMSHandler(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	if to == "" {
		http.Error(w, "Missing 'to' parameter.", http.StatusBadRequest)
		return
	}

	msg := twirest.SendMessage{
		Text: "Hello from App Engine!",
		From: twilioNumber,
		To:   to,
	}

	resp, err := twilioClient.Request(msg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not send SMS: %v", err), 500)
		return
	}

	fmt.Fprintf(w, "SMS sent successfully. Response:\n%#v", resp.Message)
}


func receiveSMSHandler(w http.ResponseWriter, r *http.Request) {
	sender := r.FormValue("From")
	body := r.FormValue("Body")

	var Url *url.URL
			Url, _ = url.Parse("http://api.wolframalpha.com/v2/result?")

			parameters := url.Values{}
			parameters.Add("appid", "WA APP ID")
			parameters.Add("i", string(body))
			Url.RawQuery = parameters.Encode()

			ctx := appengine.NewContext(r)
			client := urlfetch.Client(ctx)
	WAresp, _ := client.Get(Url.String())

	WAbody, _ := ioutil.ReadAll(WAresp.Body)

	WAfinal := (string(WAbody))

	var buffer bytes.Buffer
	buffer.WriteString(WAfinal)
	buffer.WriteString(" - W|A")


	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: fmt.Sprintf(buffer.String()),
		From: twilioNumber,
		To:   sender,
	})
	resp.Send(w)
}
