# ZUULY

## Description
Zuuly parse zuul remote config to get url and path. 

## How to use
```go
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
	RemoteConfigUrl = "http://example.com"
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
```

####Output:
```bash
2019/02/07 11:09:08 {Key:.notification}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/message}, {Front_path:/v2/message/*}

2019/02/07 11:09:08 {Key:.v3order-history}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/v3/histories}, {Front_path:/v3/histories/*}

2019/02/07 11:09:08 {Key:.registration}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/account/registration}, {Front_path:/v2/account/registration/*}

2019/02/07 11:09:08 {Key:.checkout}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/v2/checkout}, {Front_path:/v2/checkout/*}

2019/02/07 11:09:08 {Key:.tokomodal-agreement}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/tokomodal/agreement}, {Front_path:/v2/tokomodal/agreement/*}

2019/02/07 11:09:08 {Key:.tokomodal}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/tokomodal/}, {Front_path:/v2/tokomodal/*}

2019/02/07 11:09:08 {Key:.chat}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/chat}, {Front_path:/v2/chat/*}

2019/02/07 11:09:08 {Key:.kimonu}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/kimonu}, {Front_path:/v2/kimonu/*}

2019/02/07 11:09:08 {Key:.pin}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/pin}, {Front_path:/v2/auth/pin/*}

2019/02/07 11:09:08 {Key:.ppob}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/ppob}, {Front_path:/v2/account/ppob/*}

2019/02/07 11:09:08 {Key:.ppob-topup}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/ppob-topup}, {Front_path:/v2/ppob-topup/*}

2019/02/07 11:09:08 {Key:.v3checkout}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/v3/checkout}, {Front_path:/v3/checkout/*}

2019/02/07 11:09:08 {Key:.v4checkout}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/v4/checkout}, {Front_path:/v4/checkout/*}

2019/02/07 11:09:08 {Key:.tokomodal-registration}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/tokomodal/registration}, {Front_path:/tokomodal/registration/*}

2019/02/07 11:09:08 {Key:.order-history}, {Scheme:http}, {Base_url:http://example.com}, {Reverse_Proxy_Path:/v2/histories}, {Front_path:/v2/histories/*}
```
