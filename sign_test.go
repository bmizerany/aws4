package aws4

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var (
	suite []string
	sv    *Service
)

func init() {
	var err error
	suite, err = filepath.Glob(filepath.Join("aws4_testsuite", "*.req"))
	if err != nil {
		panic(err)
	}

	for i, n := range suite {
		suite[i] = n[:len(n)-len(filepath.Ext(n))]
	}

	sv = new(Service)
	sv.Name = "host"
	sv.Region = "us-east-1"
}

func readRequest(name, ext string) (string, *http.Request) {
	fr, err := os.Open(name + ".req")
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	r, err := http.ReadRequest(bufio.NewReader(fr))
	if err != nil {
		panic(err)
	}
	r.Header.Del("Content-Length")

	fe, err := os.Open(name + ext)
	if err != nil {
		panic(err)
	}
	defer fe.Close()

	b, err := ioutil.ReadAll(fe)
	if err != nil {
		panic(err)
	}

	return string(b), r
}

func TestCreateCanonicalRequest(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, name := range suite {
		exp, r := readRequest(name, ".creq")

		buf.Reset()
		sv.writeRequest(buf, r)
		if got := string(buf.Bytes()); got != exp {
			t.Logf("--- %s ---", name)
			t.Errorf("\nwanted:\n%q\ngot:\n%q", exp, got)
		}
	}
}

func TestStringToSign(t *testing.T) {
	buf := new(bytes.Buffer)

	for _, name := range suite {
		exp, r := readRequest(name, ".sts")

		buf.Reset()
		sv.writeStringToSign(buf, r)
		if got := string(buf.Bytes()); got != exp {
			t.Logf("--- %s ---", name)
			t.Errorf("\nwanted:\n%q\ngot:\n%q", exp, got)
		}
	}
}

func TestSetAuthorizationHeader(t *testing.T) {
	buf := new(bytes.Buffer)

	tk := &Keys {
		AccessKey: "AKIDEXAMPLE",
		SecretKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
	}

	for _, name := range suite {
		exp, r := readRequest(name, ".authz")
		buf.Reset()

		sv.Sign(tk, r)
		got := r.Header.Get("Authorization")
		if exp != got {
			t.Logf("--- %s ---", name)
			t.Errorf("\nwanted:\n%q\ngot:\n%q", exp, got)
		}
	}
}

func TestError(t *testing.T) {
	r, _ := http.NewRequest("POST", "http://example.com", nil)
	err := sv.Sign(new(Keys), r)
	if err != ErrNoDate {
		t.Error("expected ErrNoDate, got %#v", err)
	}
}
