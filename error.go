package petra

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type APIError struct {
	Message        string        `json:"message,omitempty"`
	HTTPStatusCode int           `json:"code,omitempty"`
	Details        ErrorResponse `json:"details,omitempty"`
	URL            *url.URL      `json:"url,omitempty"`
	Header         http.Header   `json:"header,omitempty"`
}

// APIError supports the error interface
func (aerr *APIError) Error() string {
	ret, _ := json.Marshal(aerr)
	return string(ret)
}

// ErrorResponse represents an error response from the Petra API server
type ErrorResponse struct {
	Status  bool                   `json:"status,omitempty"`
	Message string                 `json:"message,omitempty"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func newAPIError(resp *http.Response) *APIError {
	p, _ := ioutil.ReadAll(resp.Body)

	var petraErrorResp ErrorResponse
	_ = json.Unmarshal(p, &petraErrorResp)
	return &APIError{
		HTTPStatusCode: resp.StatusCode,
		Header:         resp.Header,
		Details:        petraErrorResp,
		URL:            resp.Request.URL,
	}
}
