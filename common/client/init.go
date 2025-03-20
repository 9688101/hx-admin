package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/9688101/hx-admin/core/logger"
	"github.com/9688101/hx-admin/global"
)

var HTTPClient *http.Client
var ImpatientHTTPClient *http.Client
var UserContentRequestHTTPClient *http.Client

func Init() {
	if global.UserContentRequestProxy != "" {
		logger.SysLog(fmt.Sprintf("using %s as proxy to fetch user content", global.UserContentRequestProxy))
		proxyURL, err := url.Parse(global.UserContentRequestProxy)
		if err != nil {
			logger.FatalLog(fmt.Sprintf("USER_CONTENT_REQUEST_PROXY set but invalid: %s", global.UserContentRequestProxy))
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		UserContentRequestHTTPClient = &http.Client{
			Transport: transport,
			Timeout:   time.Second * time.Duration(global.UserContentRequestTimeout),
		}
	} else {
		UserContentRequestHTTPClient = &http.Client{}
	}
	var transport http.RoundTripper
	if global.RelayProxy != "" {
		logger.SysLog(fmt.Sprintf("using %s as api relay proxy", global.RelayProxy))
		proxyURL, err := url.Parse(global.RelayProxy)
		if err != nil {
			logger.FatalLog(fmt.Sprintf("USER_CONTENT_REQUEST_PROXY set but invalid: %s", global.UserContentRequestProxy))
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	if global.RelayTimeout == 0 {
		HTTPClient = &http.Client{
			Transport: transport,
		}
	} else {
		HTTPClient = &http.Client{
			Timeout:   time.Duration(global.RelayTimeout) * time.Second,
			Transport: transport,
		}
	}

	ImpatientHTTPClient = &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
}
