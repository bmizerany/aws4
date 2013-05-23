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

var keys = &aws4.Keys{
	AccessKey: os.Getenv("AWS_ACCESS_KEY"),
	SecretKey: os.Getenv("AWS_SECRET_KEY"),
}

func init() {
	http.DefaultTransport.(*http.Transport).RegisterProtocol("aws4", &aws4.Transport{Keys: keys})
}

func Example_jSONBody() {
	data := strings.NewReader("{}")
	r, _ := http.NewRequest("POST", "aws4://dynamodb.us-east-1.amazonaws.com/", data)
	r.Header.Set("Content-Type", "application/x-amz-json-1.0")
	r.Header.Set("X-Amz-Target", "DynamoDB_20111205.ListTables")

	resp, err := http.DefaultClient.Do(r)
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

	u := "aws4://autoscaling.us-east-1.amazonaws.com/"
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	resp, err := http.Post(u, "application/x-www-form-urlencoded", body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}
