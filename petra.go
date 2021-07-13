package petra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// library version
	version = "0.1.0"

	// defaultHTTPTimeout is the default timeout on the http client
	defaultHTTPTimeout = 60 * time.Second

	// base URL for all Petra API requests
	baseURL = "https://petra-staging.herokuapp.com/"

	// User agent used when communicating with the Petra API.
	// userAgent = "Petra-go/" + version
	userAgent = "Mozilla/5.0 (Unknown; Linux) AppleWebKit/538.1 (KHTML, like Gecko) Chrome/v1.0.0 Safari/538.1"
)

type service struct {
	client *Client
}

// Client manages communication with the Petra API
type Client struct {
	common service      // Reuse a single struct instead of allocating one for each service on the heap.
	client *http.Client // HTTP client used to communicate with the API.
	// the API Key used to authenticate all Petra API requests
	key string
	baseURL *url.URL

	logger Logger
	// Services supported by the Petra API.
	// Miscellaneous actions are directly implemented on the Client object
	Customer     *CustomerService
	Transaction  *TransactionService
	Page         *PageService
	Invoice		 *InvoiceService
	LoggingEnabled bool
	Log            Logger
}

// Logger interface for custom loggers
type Logger interface {
	Printf(format string, v ...interface{})
}

// Metadata is an key-value pairs added to Petra API requests
type Metadata map[string]interface{}

// Response represents arbitrary response data
type Response map[string]interface{}

// RequestValues aliased to url.Values as a workaround
type RequestValues url.Values

// MarshalJSON to handle custom JSON decoding for RequestValues
func (v RequestValues) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, 3)
	for k, val := range v {
		m[k] = val[0]
	}
	return json.Marshal(m)
}

// ListMeta is pagination metadata for paginated responses from the Petra API
type ListMeta struct {
	Total     int `json:"total"`
	Skipped   int `json:"skipped"`
	PerPage   int `json:"perPage"`
	Page      int `json:"page"`
	PageCount int `json:"pageCount"`
}

// NewClient creates a new Petra API client with the given API key
// and HTTP client, allowing overriding of the HTTP client to use.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.
func NewClient(key string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	u, _ := url.Parse(baseURL)
	c := &Client{
		client:         httpClient,
		key:            key,
		baseURL:        u,
		LoggingEnabled: true,
		Log:            log.New(os.Stderr, "", log.LstdFlags),
	}

	c.common.client = c
	c.Customer = (*CustomerService)(&c.common)
	c.Transaction = (*TransactionService)(&c.common)
	c.Page = (*PageService)(&c.common)
	c.Invoice = (*InvoiceService)(&c.common)
	return c
}

// Call actually does the HTTP request to Petra API
func (c *Client) Call(method, path string, body, v interface{}) error {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}
	u, _ := c.baseURL.Parse(path)
	req, err := http.NewRequest(method, u.String(), buf)

	if err != nil {
		if c.LoggingEnabled {
			c.Log.Printf("Cannot create Petra request: %v\n", err)
		}
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("User-Agent", userAgent)

	if c.LoggingEnabled {
		c.Log.Printf("Requesting %v %v%v\n", req.Method, req.URL.Host, req.URL.Path)
		c.Log.Printf("POST request data %v\n", buf)
	}

	start := time.Now()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if c.LoggingEnabled {
		c.Log.Printf("Completed in %v\n", time.Since(start))
	}

	defer resp.Body.Close()
	return c.decodeResponse(resp, v)
}






// INTERNALS
func paginateURL(path string, count, offset int) string {
	return fmt.Sprintf("%s?perPage=%d&page=%d", path, count, offset)
}

func mapstruct(data interface{}, v interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result:           v,
		TagName:          "json",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	err = decoder.Decode(data)
	return err
}

func mustGetTestKey() string {
	key := os.Getenv("PETRA_KEY")

	if len(key) == 0 {
		panic("PETRA_KEY environment variable is not set\n")
	}

	return key
}

// decodeResponse decodes the JSON response from the Twitter API.
// The actual response will be written to the `v` parameter
func (c *Client) decodeResponse(httpResp *http.Response, v interface{}) error {
	var resp Response
	respBody, err := ioutil.ReadAll(httpResp.Body)
	json.Unmarshal(respBody, &resp)

	if status, _ := resp["status"].(bool); !status || httpResp.StatusCode >= 400 {
		if c.LoggingEnabled {
			c.Log.Printf("Petra error: %+v", err)
			c.Log.Printf("HTTP Response: %+v", resp)
		}
		return newAPIError(httpResp)
	}

	if c.LoggingEnabled {
		c.Log.Printf("Petra response: %v\n", resp)
	}

	if data, ok := resp["data"]; ok {
		switch t := resp["data"].(type) {
		case map[string]interface{}:
			return mapstruct(data, v)
		default:
			_ = t
			return mapstruct(resp, v)
		}
	}
	// if response data does not contain data key, map entire response to v
	return mapstruct(resp, v)
}
