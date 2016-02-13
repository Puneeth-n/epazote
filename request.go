package epazote

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const URL string = `^((ftp|https?):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`

var rxURL = regexp.MustCompile(URL)

type ServiceHttpResponse struct {
	Err     error
	Service string
}

// AsyncGet used as a URL validation method
func AsyncGet(s *Services) <-chan ServiceHttpResponse {
	ch := make(chan ServiceHttpResponse, len(*s))

	for k, v := range *s {
		go func(name string, url string, verify bool) {
			res, err := HTTPGet(url, true, v.Insecure)
			if err != nil {
				ch <- ServiceHttpResponse{err, name}
				return
			}
			res.Body.Close()
			ch <- ServiceHttpResponse{nil, name}
		}(k, v.URL, v.Insecure)
	}

	return ch
}

// HTTPGet creates a new http request
func HTTPGet(url string, follow, insecure bool, h map[string]string, timeout ...int) (*http.Response, error) {
	// timeout in seconds defaults to 5
	var t int = 5

	if len(timeout) > 0 {
		t = timeout[0]
	}

	// if insecure = true, skip ssl verification
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: insecure},
		ResponseHeaderTimeout: time.Duration(t) * time.Second,
	}

	client := &http.Client{}
	client.Transport = tr

	// create a new request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "epazote")

	// set custom headers on request
	if h != nil {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	if follow {
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// not follow redirects
	var DefaultTransport http.RoundTripper = tr

	res, err := DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// HTTPPost post service json data
func HTTPPost(url string, data []byte) error {
	// create a new request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("User-Agent", "epazote")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	res.Body.Close()

	return nil
}

// IsURL https://github.com/asaskevich/govalidator/blob/master/validator.go#L44
func IsURL(str string) bool {
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return rxURL.MatchString(str)
}

//// don't read full body
//html := io.LimitReader(resp.Body, 0)
