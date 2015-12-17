package epazote

import (
	"net/http"
	"time"
)

func Get(s Service) (*http.Response, error) {
	// timeout in seconds
	if s.Timeout == 0 {
		s.Timeout = 7
	}
	timeout := time.Duration(s.Timeout) * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	// create a new request
	req, _ := http.NewRequest("GET", s.URL, nil)
	req.Header.Set("User-Agent", "epazote")

	// try to connect
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//defer resp.Body.Close()

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
