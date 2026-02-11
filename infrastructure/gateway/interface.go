package gateway

// kontrak untuk menyimpan ke apicall
type APICall interface {
	GetAPICall() RequestAPICallResult
}

type RequestAPICallResult struct {
	RequestURL         string
	Method             string
	RequestLatency     string
	RequestBody        string
	RequestQuery       string
	ResponseBody       string
	RequestHeaders     string
	ResponseHeaders    string
	ResponseStatusCode int

	// telco tracking ID
	DotTransID  string
	ExtraField1 string
	ExtraField2 string
}
