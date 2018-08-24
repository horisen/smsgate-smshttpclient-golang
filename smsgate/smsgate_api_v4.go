package smsgate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// APIv4Preferences sets up HTTP and retry scheme
type APIv4Preferences struct {
	HTTPClient  *http.Client
	RetryPeriod time.Duration
	MaxRetries  int
	SubmitURL   string
}

// DefaultAPIv4Preferences are default API v4 preferences - default HTTP client, 3 retries after 1s.
var DefaultAPIv4Preferences = APIv4Preferences{
	HTTPClient:  http.DefaultClient,
	RetryPeriod: 1 * time.Second,
	MaxRetries:  3,
}

// APIv4Impl is API v4 protocol implementation
type APIv4Impl struct {
	preferences APIv4Preferences
}

// NewAPIv4Impl creates new instance of HTTP API with given preferences.
// If preferences are nil DefaultAPIv4Preferences will be used
func NewAPIv4Impl(preferences *APIv4Preferences) *APIv4Impl {
	apiV4 := &APIv4Impl{}

	if preferences == nil {
		apiV4.preferences = DefaultAPIv4Preferences
	} else {
		apiV4.preferences = *preferences
	}
	if apiV4.preferences.HTTPClient == nil {
		apiV4.preferences.HTTPClient = DefaultAPIv4Preferences.HTTPClient
	}
	if apiV4.preferences.RetryPeriod == 0 {
		apiV4.preferences.RetryPeriod = DefaultAPIv4Preferences.RetryPeriod
	}
	if apiV4.preferences.MaxRetries == 0 {
		apiV4.preferences.MaxRetries = DefaultAPIv4Preferences.MaxRetries
	}

	return apiV4
}

// ErrorObjWO is JSON object returned in case of error
type ErrorObjWO struct {
	Error *ErrorWO `json:"error"`
}

// ErrorWO is JSON encapsulating error report
type ErrorWO struct {
	Code        string       `json:"code"`
	Message     string       `json:"message"`
	Description string       `json:"description,omitempty"`
	Items       []*ErrorItem `json:"items,omitempty"`
}

// ErrorItem is item in ErrorWO
type ErrorItem struct {
	Name        string `json:"name"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}

// Send SMS over API.
// FIXME add retry scheme if Do returns error or API returns throughput error.
func (api *APIv4Impl) Send(req *SMSRequest) (*SMSResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reqBuf := bytes.NewBuffer(reqBytes)

	for tryNum := 0; ; tryNum++ {
		rsp, err := api.sendTry(reqBuf)
		if err == nil {
			return rsp, nil
		}
		if IsAPIError(err) {
			apiErr := err.(*APIError)
			if apiErr.Code() != RCThrottlingError {
				return rsp, err
			}
		}
		if tryNum >= api.preferences.MaxRetries {
			return rsp, err
		}
		time.Sleep(api.preferences.RetryPeriod)
	}
}

func (api *APIv4Impl) sendTry(reqBuf *bytes.Buffer) (*SMSResponse, error) {
	log.Printf("Sending SMS to %s\n", api.preferences.SubmitURL)
	httpRequest, err := http.NewRequest("POST", api.preferences.SubmitURL, reqBuf)
	if err != nil {
		return nil, err
	}

	httpClient := api.preferences.HTTPClient

	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 && resp.StatusCode != 420 {
		return nil, fmt.Errorf("Returned error code %d instead of 202 or 420", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 202 {
		var respJSON SMSResponse
		if err := json.Unmarshal(respBytes, &respJSON); err != nil {
			return nil, err
		}
		return &respJSON, nil
	}
	var errJSON ErrorObjWO
	if err := json.Unmarshal(respBytes, &errJSON); err != nil {
		return nil, err
	}

	if errJSON.Error == nil {
		return nil, NewAPIError("Unknown API Error", RCApplicationError)
	}
	errCode, err := strconv.ParseInt(errJSON.Error.Code, 10, 16)
	if err != nil {
		return nil, NewAPIError("Error parsing API Error Code", RCApplicationError)
	}

	return nil, NewAPIError(errJSON.Error.Message, ErrorCode(errCode))
}

// ParseDeliveryReport should be used to transform received HTTP request of DLR into JSONs,
// with DeliveryReport or error
func (api *APIv4Impl) ParseDeliveryReport(req *http.Request) (*DeliveryReport, error) {
	reqBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var dlrJSON DeliveryReport
	if err := json.Unmarshal(reqBytes, &dlrJSON); err != nil {
		return nil, err
	}
	return &dlrJSON, nil
}
