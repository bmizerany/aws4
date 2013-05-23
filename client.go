package aws4

import (
	"io"
	"net/http"
)

// Post is like net/http.Post but signs the request with keys before sending.
func Post(url, bodyType string, body io.Reader, keys *Keys) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return Do(req, keys)
}

// Do is like net/http.Do but signs the request with keys before sending.
func Do(req *http.Request, keys *Keys) (*http.Response, error) {
	Sign(keys, req)
	return http.DefaultClient.Do(req)
}
