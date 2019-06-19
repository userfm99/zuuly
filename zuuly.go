package zuuly

import (
	"encoding/json"
	"fmt"
	"git.alfacart.com/mobile-api/zuuly/httpclient"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type FilterKeyFunc func(key *string) *string

type FilterFunc func(field *string) (*string, bool)

type Skip func(key *string) (bool, bool)

type Config struct {
	Skip         Skip
	BaseRouteUrl []string
}

type CloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}

type Proxy struct {
	Routes map[string]ZuulRoute
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

type ZuulRoute struct {
	ReverseProxyScheme  string
	ReverseProxyBaseURL string
	ReverseProxyPath    string
	FrontPath           string
}

type RemoteConfig struct {
	URL struct {
		AccountService     string `valid:"required"`
		ProductService     string `valid:"required"`
		PromotionService   string `valid:"required"`
		TransactionService string `valid:"required"`
		TokomodalService   string `valid:"required"`
		OptionsService     string `valid:"required"`
	}
	JWT struct {
		ExpiredAtSeconds int64 `valid:"required"`
	}
	ZuulRoutes map[string]ZuulRoute
}

type (
	Zuuly struct {
		RequestAttr *httpclient.RequestAttr
	}
)

func New(RequestAttr *httpclient.RequestAttr) *Zuuly {
	return &Zuuly{RequestAttr: RequestAttr}
}

func (z *Zuuly) GetProxy(filterKeyFunc FilterKeyFunc) (*Proxy, error) {
	c := httpclient.New(httpclient.DefaultClientTimeOut)
	resp, err := c.Exchange(z.RequestAttr)
	if err != nil {
		return nil, httpclient.NewHTTPError(http.StatusInternalServerError, "ZUULY: "+err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, httpclient.NewHTTPError(http.StatusBadRequest, "ZUULY: "+err.Error())
	}
	defer resp.Body.Close()

	var cloudConfig CloudConfig
	b, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &cloudConfig)
	if err != nil {
		return nil, err
	}

	// Parse zuul routes
	zuulRoutes := make(map[string]ZuulRoute)
	propertySources := cloudConfig.PropertySources[0].Source
	var key *string
	for k, v := range propertySources {
		if key = &k; filterKeyFunc != nil {
			key = filterKeyFunc(&k)
		}
		if key != nil {
			if strings.HasSuffix(*key, ".url") {
				*key = strings.Replace(*key, ".url", "", 1)
				u, _ := url.Parse(zuulURL(v.(string), propertySources))
				if zr, ok := zuulRoutes[*key]; ok {
					zr.ReverseProxyBaseURL = u.Host
					zr.ReverseProxyScheme = u.Scheme
					zr.ReverseProxyPath = u.Path
					zuulRoutes[*key] = zr
				} else {
					zuulRoutes[*key] = ZuulRoute{
						ReverseProxyBaseURL: u.Host,
						ReverseProxyScheme:  u.Scheme,
						ReverseProxyPath:    u.Path,
					}
				}
			} else if strings.HasSuffix(k, ".path") {
				*key = strings.Replace(*key, ".path", "", 1)
				strVal := v.(string)
				strVal = strings.Replace(strVal, "/**", "/*", 1)
				if zr, ok := zuulRoutes[*key]; ok {
					zr.FrontPath = strVal
					zuulRoutes[*key] = zr
				} else {
					zuulRoutes[*key] = ZuulRoute{
						FrontPath: strVal,
					}
				}
			}
		}
	}
	return &Proxy{Routes: zuulRoutes}, nil
}

func (z *Zuuly) GetKey(filter FilterFunc) FilterKeyFunc {
	return func(key *string) *string {
		if v, ok := filter(key); ok {
			return v
		}
		return nil
	}
}

// when spring cloud config value is ${other.key}, then find the value of "other.key"
// then parse the value of other.key, wether it's OS env value of the value of other.key itself
// e.g.:
// 		 env value of SOME_OTHER_ENV_KEY=http://google.com
//       other.key.url: "${SOME_OTHER_ENV_KEY}"
//       url: ${other.key.url}/api/path
// thus, value of url = http://google.com/api/path
func zuulURL(route string, props map[string]interface{}) string {
	reg, err := regexp.Compile(`\$\{(.*?)\}`)
	if err != nil {
		return route
	}
	match := reg.FindStringSubmatch(route)
	if len(match) != 2 {
		return route
	}

	// found the key
	key, ok := props[match[1]]
	if !ok {
		return ""
	}

	if keyStr, ok := key.(string); ok {
		hostPath := strings.Split(route, "}")
		return fmt.Sprintf("%s%s", envValue(keyStr), hostPath[1])
	}
	return ""
}

func envValue(envStringKey string) string {
	reg, err := regexp.Compile(`\$\{(.*?)\}`)
	if err != nil {
		return envStringKey
	}
	match := reg.FindStringSubmatch(envStringKey)
	if len(match) != 2 {
		return envStringKey
	}
	return os.Getenv(match[1])
}
