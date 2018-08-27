# GO library for SMS Gate

## Installation

Install GO library

```
go get github.com/horisen/smsgate-clientlib-golang/smsgate
```

## Send SMS

```golang
import "github.com/horisen/smsgate-clientlib-golang/smsgate"

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
}
```