package eureka

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func executeQuery(httpAction HttpAction) ([]byte, error) {
	req := buildHttpRequest(httpAction)

	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := DefaultTransport.RoundTrip(req)
	if err != nil {
		return []byte(nil), err
	} else {
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte(nil), err
		}
		return responseBody, nil
	}
}

func doHttpRequest(httpAction HttpAction) bool {
	req := buildHttpRequest(httpAction)

	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := DefaultTransport.RoundTrip(req)
	if resp != nil && resp.StatusCode > 299 {
		defer resp.Body.Close()
		log.Printf("HTTP request failed with status: %d", resp.StatusCode)
		return false
	} else if err != nil {
		log.Printf("HTTP request failed: %s", err.Error())
		return false
	} else {
		return true
		defer resp.Body.Close()
	}
	return false
}

func buildHttpRequest(httpAction HttpAction) *http.Request {
	var req *http.Request
	var err error
	if httpAction.Body != "" {
		reader := strings.NewReader(httpAction.Body)
		req, err = http.NewRequest(httpAction.Method, httpAction.Url, reader)
	} else if httpAction.Template != "" {
		reader := strings.NewReader(httpAction.Template)
		req, err = http.NewRequest(httpAction.Method, httpAction.Url, reader)
	} else {
		req, err = http.NewRequest(httpAction.Method, httpAction.Url, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Add headers
	req.Header = map[string][]string{
		"Accept":       {httpAction.Accept},
		"Content-Type": {httpAction.ContentType},
	}

	if strings.Contains(httpAction.Url, "@") {
		start := strings.Index(httpAction.Url, "http://")
		end := strings.Index(httpAction.Url, "@")
		up := httpAction.Url[start+7 : end]
		ups := strings.Split(up, ":")
		if len(ups) == 2 {
			user := ups[0]
			passwd := ups[1]
			req.SetBasicAuth(user, passwd)
		}
	}
	return req
}
