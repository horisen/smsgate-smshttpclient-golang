# GO library for SMS Gate

## Installation

Install GO library

```
go get github.com/horisen/smsgate-smshttpclient-golang/smsgate
```

## Send SMS

```golang
import "github.com/horisen/smsgate-smshttpclient-golang/smsgate"

...

api := smsgate.NewAPIv4Impl(
    &smsgate.APIv4Preferences{
        SubmitURL: "https://SMSGATE-SUBMISSION-URL",
    })
sms := &smsgate.SMSRequest{
    Type:     smsgate.TypeText,
    Sender:   "SMS-SENDER",
    Receiver: "SMS-RECEIVER-MSISDN",

    Text: paramText,
    DCS:  smsgate.DCSGSM,
    DlrMask: smsgate.DLRMaskStandard,
    DlrURL:  "http://YOUR-SERVER-IP:YOUR-SERVER-PORT/dlr",

     Auth: &smsgate.Auth{
        Username: "YOUR-USERNAME",
        Password: "YOUR-PASSWORD",
    },
}

response, err := api.Send(sms)
if err != nil {
    if smsgate.IsAPIError(err) {
        apiErr := err.(*smsgate.APIError)
        log.Errorf("API returned error: %d : %s\n",
            apiErr.Code(), apiErr.Error())
    } else {
        log.Errorf("Error: %s\n", err)
    }
} else {
    log.Printf("Sent as message ID %s with %d parts\n",
		response.MsgID, response.NumParts)
}
```

## Receive DLRs

```golang
	mux := http.NewServeMux()
	mux.HandleFunc("/dlr", func(w http.ResponseWriter, req *http.Request) {
		dlr, err := api.ParseDeliveryReport(req)
		if err != nil {
			log.Printf("Cannot parse DLR: %s\n", err)
		}
		log.Printf("Received DLR: %#v\n", dlr)
	})

	httpsrv := &http.Server{
		Addr:           ":YOUR-SERVER-PORT",
		Handler:        mux,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(httpsrv.ListenAndServe())
```

Check `examples` directory.
