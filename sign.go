package aws4

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"runtime"
)

var ErrNoDate = errors.New("X-Amz-Date or Date header not supplied")

var crlf = []byte{'\n'}

type Keys struct {
	AccessKey string
	SecretKey string
}

func (k *Keys) sign(s *Service, t time.Time) []byte {
	h := ghmac([]byte("AWS4"+k.SecretKey), []byte(t.Format("20060102")))
	h = ghmac(h, []byte(s.Region))
	h = ghmac(h, []byte(s.Name))
	h = ghmac(h, []byte("aws4_request"))
	return h
}

type Service struct {
	Name   string
	Region string
}

func (s *Service) Sign(keys *Keys, r *http.Request) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			switch v := e.(type) {
			case nil:
			case runtime.Error:
				panic(v)
			case error:
				err = v
			default:
				panic(v)
			}
		}
	}()

	t := getDate(r)
	k := keys.sign(s, t)
	h := hmac.New(sha256.New, k)
	s.writeStringToSign(h, r)

	auth := bytes.NewBufferString("AWS4-HMAC-SHA256 ")
	auth.Write([]byte("Credential=" + keys.AccessKey + "/" + s.creds(t)))
	auth.Write([]byte{',', ' '})
	auth.Write([]byte("SignedHeaders="))
	s.writeHeaderList(auth, r)
	auth.Write([]byte{',', ' '})
	auth.Write([]byte("Signature=" + fmt.Sprintf("%x", h.Sum(nil))))

	r.Header.Set("Authorization", auth.String())

	return
}

func (s *Service) writeQuery(w io.Writer, r *http.Request) {
	var a []string
	for k, vs := range r.URL.Query() {
		k = url.QueryEscape(k)
		for _, v := range vs {
			if v == "" {
				a = append(a, k)
			} else {
				v = url.QueryEscape(v)
				a = append(a, k+"="+v)
			}
		}
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write([]byte{'&'})
		}
		w.Write([]byte(s))
	}
}

func (s *Service) writeHeader(w io.Writer, r *http.Request) {
	i, a := 0, make([]string, len(r.Header))
	for k, v := range r.Header {
		sort.Strings(v)
		a[i] = strings.ToLower(k) + ":" + strings.Join(v, ",")
		i++
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write(crlf)
		}
		io.WriteString(w, s)
	}
}

func (s *Service) writeHeaderList(w io.Writer, r *http.Request) {
	i, a := 0, make([]string, len(r.Header))
	for k, _ := range r.Header {
		a[i] = strings.ToLower(k)
		i++
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write([]byte{';'})
		}
		w.Write([]byte(s))
	}
}

func (s *Service) writeBody(w io.Writer, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	h := sha256.New()
	h.Write(b)
	fmt.Fprintf(w, "%x", h.Sum(nil))
}

func (s *Service) writeURI(w io.Writer, r *http.Request) {
	path := r.URL.RequestURI()
	if r.URL.RawQuery != "" {
		path = path[:len(path)-len(r.URL.RawQuery)-1]
	}
	slash := strings.HasSuffix(path, "/")
	path = filepath.Clean(path)
	if path != "/" && slash {
		path += "/"
	}
	w.Write([]byte(path))
}

func (s *Service) writeRequest(w io.Writer, r *http.Request) {
	r.Header.Set("host", r.Host)

	w.Write([]byte(r.Method))
	w.Write(crlf)
	s.writeURI(w, r)
	w.Write(crlf)
	s.writeQuery(w, r)
	w.Write(crlf)
	s.writeHeader(w, r)
	w.Write(crlf)
	w.Write(crlf)
	s.writeHeaderList(w, r)
	w.Write(crlf)
	s.writeBody(w, r)
}

func (s *Service) writeStringToSign(w io.Writer, r *http.Request) {
	t := getDate(r)

	w.Write([]byte("AWS4-HMAC-SHA256"))
	w.Write(crlf)
	w.Write([]byte(t.Format("20060102T150405Z")))
	w.Write(crlf)

	w.Write([]byte(s.creds(t)))
	w.Write(crlf)

	h := sha256.New()
	s.writeRequest(h, r)
	fmt.Fprintf(w, "%x", h.Sum(nil))
}

func (s *Service) creds(t time.Time) string {
	return t.Format("20060102") + "/" + s.Region + "/" + s.Name + "/aws4_request"
}

func ghmac(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func getDate(r *http.Request) time.Time {
	date := r.Header.Get("X-Amz-Date")
	if date == "" {
		date = r.Header.Get("Date")
		if date == "" {
			panic(ErrNoDate)
		}
	}
	t, err := time.Parse(http.TimeFormat, date)
	if err != nil {
		panic(err)
	}
	return t
}
