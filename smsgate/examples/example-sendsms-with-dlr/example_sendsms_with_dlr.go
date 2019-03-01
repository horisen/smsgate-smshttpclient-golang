package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/horisen/smsgate-smshttpclient-golang/smsgate"
)

const (
	// DefaultDLRServerPort is default port  where we listen to DLR callbacks
	DefaultDLRServerPort = 17777
)

var (
	paramSender        string
	paramReceiver      string
	paramText          string
	paramUsername      string
	paramPassword      string
	paramSubmitURL     string
	paramDLRServerHost string
	paramDLRServerPort int
)

func init() {
	flag.StringVar(&paramSender, "sender", "", "SMS Sender")
	flag.StringVar(&paramReceiver, "receiver", "", "SMS Receiver")
	flag.StringVar(&paramText, "text", "", "SMS Text")
	flag.StringVar(&paramUsername, "username", "", "Username")
	flag.StringVar(&paramPassword, "password", "", "SMS Password")
	flag.StringVar(&paramSubmitURL, "submiturl", "", "Submit URL")
	flag.StringVar(&paramDLRServerHost, "dlrhost", "", "DLR Server Host")
	flag.IntVar(&paramDLRServerPort, "dlrport", DefaultDLRServerPort, "DLR Server Port")
}

func main() {
	flag.Parse()
	api := smsgate.NewAPIv4Impl(
		&smsgate.APIv4Preferences{
			SubmitURL: paramSubmitURL,
		})

	dlrCaught := make(chan *smsgate.DeliveryReport)

	dlrServer := startDLRServer(api, dlrCaught)

	sms := &smsgate.SMSRequest{
		Type:     smsgate.TypeText,
		Sender:   paramSender,
		Receiver: paramReceiver,

		Text:    paramText,
		DCS:     smsgate.DCSGSM,
		DlrMask: smsgate.DLRMaskStandard,
		DlrURL:  fmt.Sprintf("http://%s:%d/dlr", paramDLRServerHost, paramDLRServerPort),

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

	select {
	case dlr := <-dlrCaught:
		log.Printf("Caught DLR: %#v\n", dlr)
	case <-time.After(1 * time.Minute):
		log.Printf("DLR didn't arrive after 1min\n")
	}

	dlrServer.Shutdown(context.Background())
}

func startDLRServer(api smsgate.API, dlrChan chan *smsgate.DeliveryReport) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/dlr", func(w http.ResponseWriter, req *http.Request) {
		dlr, err := api.ParseDeliveryReport(req)
		if err != nil {
			log.Printf("Cannot parse DLR: %s\n", err)
		}
		log.Printf("Received DLR: %#v\n", dlr)
		dlrChan <- dlr
	})

	httpsrv := &http.Server{
		Addr:           fmt.Sprintf(":%d", paramDLRServerPort),
		Handler:        mux,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		log.Fatal(httpsrv.ListenAndServe())
	}()
	return httpsrv
}
