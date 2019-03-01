package main

import (
	"flag"
	"log"

	"github.com/horisen/smsgate-smshttpclient-golang/smsgate"
)

var (
	paramSender    string
	paramReceiver  string
	paramText      string
	paramUsername  string
	paramPassword  string
	paramSubmitURL string
)

func init() {
	flag.StringVar(&paramSender, "sender", "", "SMS Sender")
	flag.StringVar(&paramReceiver, "receiver", "", "SMS Receiver")
	flag.StringVar(&paramText, "text", "", "SMS Text")
	flag.StringVar(&paramUsername, "username", "", "Username")
	flag.StringVar(&paramPassword, "password", "", "SMS Password")
	flag.StringVar(&paramSubmitURL, "submiturl", "", "Submit URL")
}

func main() {
	flag.Parse()

	api := smsgate.NewAPIv4Impl(
		&smsgate.APIv4Preferences{
			SubmitURL: paramSubmitURL,
		})
	sms := &smsgate.SMSRequest{
		Type:     smsgate.TypeText,
		Sender:   paramSender,
		Receiver: paramReceiver,

		Text: paramText,
		DCS:  smsgate.DCSGSM,

		Auth: &smsgate.Auth{
			Username: paramUsername,
			Password: paramPassword,
		},
	}
	response, err := api.Send(sms)
	if err != nil {
		if smsgate.IsAPIError(err) {
			apiErr := err.(*smsgate.APIError)
			log.Fatalf("API returned error: %d : %s\n",
				apiErr.Code(), apiErr.Error())
		}
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Sent as message ID %s with %d parts\n",
		response.MsgID, response.NumParts)
}
