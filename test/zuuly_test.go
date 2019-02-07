package test

import (
	"alfacart/zuuly"
	"alfacart/zuuly/httpclient"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var (
	Prefix          = "zuul.routes"
	RemoteConfigUrl = "http://10.222.12.46:7000/alfamikro-mapi-gateway/stage"
)

func TestZuuly(t *testing.T) {
	u, err := url.Parse(RemoteConfigUrl)
	if err != nil {
		t.Fatalf("error parsing url %s", err.Error())
	}
	z := zuuly.New(&httpclient.RequestAttr{Url: u, Method: http.MethodGet})
	proxy, err := z.GetProxy(z.GetKey(func(field *string) (*string, bool) {
		if strings.HasPrefix(*field, Prefix) {
			s := strings.Replace(*field, Prefix, "", 1)
			return &s, true
		}
		return nil, false
	}))
	if err != nil {
		t.Error(err.Error())
	} else {
		for k, v := range proxy.Routes {
			log.Printf("{Key:%s}, {Scheme:%s}, {Base_url:%s}, {Reverse_Proxy_Path:%s}, {Front_path:%s}\n", k, v.ReverseProxyScheme, v.ReverseProxyBaseURL, v.ReverseProxyPath, v.FrontPath)
		}
	}
}
