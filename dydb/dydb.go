// This is an experimental library for use with DynamoDB. It uses
// github.com/bmizerany/aws4 to sign requests. See Example for use.
package dydb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bmizerany/aws4"
	"io"
	"net/http"
)

// A ResponseError is returned by Decode when an error communicating with
// DynamoDB occurs.
type ResponseError struct {
	StatusCode int
	Body       io.Reader
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("dydb: response error: %d", e.StatusCode)
}

type errorDecoder struct {
	err error
}

func (e *errorDecoder) Decode(v interface{}) error {
	return e.err
}

type Decoder interface {
	Decode(v interface{}) error
}

const DefaultURL = "https://dynamodb.us-east-1.amazonaws.com/"

type DB struct {
	// The version of DynamoDB to use (default is latest)
	Version string

	// If nil, aws4.DefaultClient is used.
	Client *aws4.Client

	// If empty, DefaultURL is used.
	URL string
}

// Do executes action with and JSON-encoded v as the body. If v is nil, an empty {} is used as the body.
func (c *DB) Do(action string, v interface{}) Decoder {
	cl := c.Client
	if cl == nil {
		cl = aws4.DefaultClient
	}

	url := c.URL
	if url == "" {
		url = DefaultURL
	}

	ver := c.Version
	if ver == "" {
		ver = "20120810"
	}

	if v == nil {
		v = struct{}{}
	}

	b, err := json.Marshal(v)
	if err != nil {
		return &errorDecoder{err: err}
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return &errorDecoder{err: err}
	}
	r.Header.Set("Content-Type", "application/x-amz-json-1.0")
	r.Header.Set("X-Amz-Target", "DynamoDB_"+ver+"."+action)

	resp, err := cl.Do(r)
	if err != nil {
		return &errorDecoder{err: err}
	}

	if code := resp.StatusCode; code != 200 {
		// Read the whole body in so that Keep-Alives may be released back to the pool.
		b := new(bytes.Buffer)
		io.Copy(b, resp.Body)
		return &errorDecoder{err: &ResponseError{StatusCode: code, Body: b}}
	}
	return json.NewDecoder(resp.Body)
}
