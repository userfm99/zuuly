package httpclient

import (
	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"

	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBackoffInterval   = 2 * time.Millisecond
	DefaultMaxJitterInterval = 5 * time.Millisecond
	DefaultClientTimeOut     = 3000 * time.Millisecond
	DefaultRetryCount        = 4
)

type RequestAttr struct {
	Url       *url.URL // Required
	Method    string   // Required
	Body      io.Reader
	HeaderMap map[string]string
	Timeout   time.Duration
}

type C struct {
	TimeOut time.Duration
}

func New(timeOut time.Duration) *C {
	var c = C{TimeOut: DefaultClientTimeOut}
	if timeOut >= DefaultClientTimeOut {
		c.TimeOut = timeOut
	}
	return &c
}

func (c *C) Exchange(requestAttr *RequestAttr) (*http.Response, error) {
	if err := validateRequestAttr(requestAttr); err != nil {
		return nil, err
	}

	hc := httpclient.NewClient(
		httpclient.WithHTTPTimeout(c.TimeOut),
		httpclient.WithRetrier(Retrier()),
		httpclient.WithRetryCount(DefaultRetryCount))

	req, err := http.NewRequest(requestAttr.Method, requestAttr.Url.String(), requestAttr.Body)
	if err != nil {
		return nil, err
	}

	return hc.Do(req)
}
func validateRequestAttr(attr *RequestAttr) error {
	if attr.Url == nil || attr.Method == "" {
		return NewHTTPError(http.StatusBadRequest, "URL or method not specified")
	}
	return nil
}

func Retrier() heimdall.Retriable {
	return heimdall.NewRetrier(DefaultBackoff(DefaultBackoffInterval, DefaultMaxJitterInterval))
}
func DefaultBackoff(backoffInterval time.Duration, jitterInterval time.Duration) heimdall.Backoff {
	return heimdall.NewConstantBackoff(backoffInterval, jitterInterval)
}
