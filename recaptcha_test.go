package recaptcha

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// https://developers.google.com/recaptcha/docs/faq#id-like-to-run-automated-tests-with-recaptcha.-what-should-i-do
const (
	testSecret = "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe"
	// testSitekey = "6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI"
)

func TestVerifyV2(t *testing.T) {
	testCases := []struct {
		description    string
		clientResponse string
		secret         string
		err            error
	}{
		{
			description:    "happy",
			secret:         testSecret, // secret always leads to success in V2
			clientResponse: "anything",
		},
		{
			description: "missing form",
			err:         ErrNoCaptcha,
		},
		{
			description:    "missing secret",
			clientResponse: "anything",
			err:            ErrNoSuccess,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			req := newRequest(t, tt.clientResponse)
			_, err := VerifyV2(tt.secret, req)
			if !errors.Is(err, tt.err) {
				t.Fatalf("want error '%v' but got '%v'", tt.err, err)
			}
		})
	}
}

func TestVerifyV3(t *testing.T) {
	testCases := []struct {
		description    string
		clientResponse string
		secret         string
		options        []OptionV3
		err            error
	}{
		{
			description:    "happy",
			secret:         testSecret, // secret always leads to score 0.0 in V3
			clientResponse: "anything",
			options:        []OptionV3{MinScore(0.0)},
		},
		{
			description: "missing_form",
			err:         ErrNoCaptcha,
		},
		{
			description:    "missing_secret",
			clientResponse: "wrong",
			err:            ErrNoSuccess,
		},
		{
			description:    "no_success",
			secret:         testSecret,
			clientResponse: "anything",
			err:            ErrScore,
		},

		{
			description:    "wrong_action",
			secret:         testSecret,
			clientResponse: "anything",
			options:        []OptionV3{MinScore(0.0), Action("register")},
			err:            ErrAction,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			req := newRequest(t, tt.clientResponse)
			_, err := VerifyV3(tt.secret, req, tt.options...)
			if !errors.Is(err, tt.err) {
				t.Fatalf("want error '%v' but got '%v'", tt.err, err)
			}
		})
	}
}

func newRequest(t *testing.T, clientResponse string) *http.Request {
	form := url.Values{}
	form.Set("g-recaptcha-response", clientResponse)

	req, err := http.NewRequest(http.MethodPost, "https://example.com", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.RemoteAddr = "127.0.0.1:58662"
	return req
}
