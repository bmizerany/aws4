package aws4_test

import (
	"fmt"
	"github.com/bmizerany/aws4"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func Example_jSONBody() {
	data := strings.NewReader("{}")
	r, _ := http.NewRequest("POST", "https://dynamodb.us-east-1.amazonaws.com/", data)
	r.Header.Set("Host", "dynamodb.us-east-1.awsamazon.com")
	r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	r.Header.Set("Content-Type", "application/x-amz-json-1.0")
	r.Header.Set("X-Amz-Target", "DynamoDB_20111205.ListTables")

	tk := &aws4.Keys{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
	}

	sv := &aws4.Service{
		Name:   "dynamodb",
		Region: "us-east-1",
	}

	if err := sv.Sign(tk, r); err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}

func Example_formEncodedBody() {
	r, _ := http.NewRequest("POST", "https://autoscaling.us-east-1.amazonaws.com/", nil)
	r.Header.Set("Host", "autoscaling.us-east-1.amazonaws.com")
	r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	v := make(url.Values)
	v.Set("Action", "DescribeAutoScalingGroups")

	r.Body = ioutil.NopCloser(strings.NewReader(v.Encode()))

	tk := &aws4.Keys{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
	}

	sv := &aws4.Service{
		Name:   "autoscaling",
		Region: "us-east-1",
	}

	if err := sv.Sign(tk, r); err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}
