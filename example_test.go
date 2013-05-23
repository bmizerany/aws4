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
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var keys = &aws4.Keys{
	AccessKey: os.Getenv("AWS_ACCESS_KEY"),
	SecretKey: os.Getenv("AWS_SECRET_KEY"),
}

func Example_jSONBody() {
	data := strings.NewReader("{}")
	r, _ := http.NewRequest("POST", "https://dynamodb.us-east-1.amazonaws.com/", data)
	r.Header.Set("Content-Type", "application/x-amz-json-1.0")
	r.Header.Set("X-Amz-Target", "DynamoDB_20111205.ListTables")

	resp, err := aws4.Do(r, keys)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}

func Example_formEncodedBody() {
	v := make(url.Values)
	v.Set("Action", "DescribeAutoScalingGroups")

	url := "https://autoscaling.us-east-1.amazonaws.com/"
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	resp, err := aws4.Post(url, "application/x-www-form-urlencoded", body, keys)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}
