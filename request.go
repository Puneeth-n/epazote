package epazote

import (
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
func AsyncGet(s Services) <-chan ServiceHttpResponse {
	ch := make(chan ServiceHttpResponse, len(s))

	for k, v := range s {
		go func(name string, url string) {
			resp, err := Get(url)
			if err != nil {
				ch <- ServiceHttpResponse{err, name}
				return
			}
			resp.Body.Close()
			ch <- ServiceHttpResponse{nil, name}
		}(k, v.URL)
	}

	return ch
}

// Get creates a new http request
func Get(url string, timeout ...int) (*http.Response, error) {
	// timeout in seconds defaults to 5
	var t int = 5

	if len(timeout) > 0 {
		t = timeout[0]
	}

	client := &http.Client{
		Timeout: time.Duration(t) * time.Second,
	}

	// create a new request
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "epazote")

	// try to connect
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

//if len(s.Expect.Body) > 0 {
//// read full body
//html = resp.Body
//}

//// read the body
//body, err := ioutil.ReadAll(html)
//if err != nil {
//return err
//}
