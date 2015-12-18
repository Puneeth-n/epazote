package epazote

import (
	"net/http"
	"time"
)

type ServiceHttpResponse struct {
	Err     error
	Service string
}

func AsyncGet(services map[string]Service) <-chan ServiceHttpResponse {
	ch := make(chan ServiceHttpResponse, len(services))

	for k, v := range services {
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
