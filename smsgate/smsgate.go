// Package smsgate provides wrapper for accessing SMS BulkGate service
package smsgate

// CustomParameters is type for storing custom parameters in request (JSON of any structure)
type CustomParameters map[string]interface{}

// Auth is structure for storing authentication data in SMS request
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ContentType is type of SMS content
type ContentType string

const (
	// TypeText is SMS type for sending text SMS
	TypeText ContentType = "text"
	// TypeWSI is SMS type for sending Binary WAP Service Indication SMS
	TypeWSI ContentType = "wsi"
)

// DataCodingScheme is coding scheme for SMS
type DataCodingScheme string

const (
	// DCSGSM is data coding scheme - GSM
	DCSGSM = "gsm"
	// DCSUCS is data coding scheme - UCS2
	DCSUCS = "ucs"
)

// DeliveryReportMask is mask of events we want to receive
type DeliveryReportMask int

// Delivery report mask values
const (
	DLRMaskDelivered   DeliveryReportMask = 1
	DLRMaskUndelivered DeliveryReportMask = 2
	DLRMaskBuffered    DeliveryReportMask = 4
	DLRMaskSentToSMSC  DeliveryReportMask = 8
	DLRMaskRejected    DeliveryReportMask = 16

	DLRMaskNone DeliveryReportMask = 0

	DLRMaskStandard = (DLRMaskDelivered | DLRMaskUndelivered | DLRMaskRejected)
)

// SMSRequest is structure for submitting SMS
type SMSRequest struct {
	// Common for all types
	Type     ContentType        `json:"type"`
	Sender   string             `json:"sender"`
	Receiver string             `json:"receiver"`
	DlrMask  DeliveryReportMask `json:"dlrMask"`
	DlrURL   string             `json:"dlrUrl"`
	Flash    bool               `json:"flash"`

	// Type=Text
	Text string           `json:"text"`
	DCS  DataCodingScheme `json:"dcs"`

	// Type=WSI
	URL   string `json:"url"`
	Title string `json:"title"`

	Custom CustomParameters `json:"custom,omitempty"`

	Auth *Auth `json:"auth"`
}

// SMSResponse is structure used for SMS Response JSON Encoding.
type SMSResponse struct {
	MsgID    string `json:"msgId,omitempty"`
	NumParts int    `json:"numParts,omitempty"`
}

const (
	// DLREventDelivered is DLR event when SMS is delivered
	DLREventDelivered = "DELIVERED"
	// DLREventUndelivered is DLR event when SMS is udelivered
	DLREventUndelivered = "UNDELIVERED"
	// DLREventBuffered  is DLR event when SMS is buffered
	DLREventBuffered = "BUFFERED"
	// DLREventSentToSMSC  is DLR event when SMS is sent to SMSC
	DLREventSentToSMSC = "SENT_TO_SMSC"
	// DLREventRejected  is DLR event when SMS is rejected
	DLREventRejected = "REJECTED"
	// DLREventExpired  is DLR event when SMS is expired
	DLREventExpired = "EXPIRED"
	// DLREventUnknown  is DLR event when SMS is unknown
	DLREventUnknown = "UNKNOWN"
)

// DeliveryReport os structure representing delivery report
type DeliveryReport struct {
	MsgID        string `json:"msgId"`
	Event        string `json:"event"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	PartNum      int    `json:"partNum"`
	NumParts     int    `json:"numParts"`
	AccountName  string `json:"accountName"`
	SendTime     int    `json:"sendTime"`
	DlrTime      int    `json:"dlrTime"`

	Custom CustomParameters `json:"custom,omitempty"`
}
