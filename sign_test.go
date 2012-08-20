package aws4

import (
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	r, _ := http.NewRequest("POST", "http://example.com", nil)
	sv := &Service{}
	err := sv.Sign(new(Keys), r)
	if err != ErrNoDate {
		t.Error("expected ErrNoDate, got %#v", err)
	}
}
