// based on https://www.socketloop.com/tutorials/golang-recaptcha-example
package recaptcha

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// ResponseV2 is a reCAPTCHA v2 response.
type ResponseV2 struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"` // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname    string    `json:"hostname"`     // the hostname of the site where the reCAPTCHA was solved
	ErrorCodes  []string  `json:"error-codes"`  // optional
}

// VerifyV2 verifies a reCAPTCHA v2 request.
func VerifyV2(secret string, r *http.Request) (*ResponseV2, error) {
	body, err := verify(secret, r)
	if err != nil {
		return nil, err
	}

	var resp ResponseV2
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return &resp, ErrNoSuccess
	}
	return &resp, nil
}

// ResponseV3 is a reCAPTCHA v3 response.
type ResponseV3 struct {
	Score  float64 `json:"score"`
	Action string  `json:"action"`
	ResponseV2
}

type optionsV3 struct {
	minScore float64
	action   string
}

// OptionV3 are optional settings for reCAPTCHA V3 validations.
type OptionV3 func(*optionsV3)

// MinScore sets the minimum required score for a successful validation (default: 0.5).
func MinScore(min float64) OptionV3 {
	return func(f *optionsV3) {
		f.minScore = min
	}
}

// Action sets the required action for a successful validation (default is empty).
func Action(action string) OptionV3 {
	return func(f *optionsV3) {
		f.action = action
	}
}

// VerifyV3 verifies a reCAPTCHA v3 request.
func VerifyV3(secret string, r *http.Request, opts ...OptionV3) (*ResponseV3, error) {
	// Default options
	o := &optionsV3{
		minScore: 0.5,
	}

	for _, opt := range opts {
		opt(o)
	}

	body, err := verify(secret, r)
	if err != nil {
		return nil, err
	}

	var resp ResponseV3
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return &resp, ErrNoSuccess
	}
	if resp.Score < o.minScore {
		return &resp, fmt.Errorf("%w: %v/%v", ErrScore, resp.Score, o.minScore)
	}
	if o.action != "" && resp.Action != o.action {
		return &resp, fmt.Errorf("%w: want '%s' but got '%s'", ErrAction, o.action, resp.Action)
	}

	return &resp, nil
}

var (
	// ErrNoCaptcha is returned when the form value 'g-recaptcha-response' is empty
	ErrNoCaptcha = errors.New("missing recaptcha response in request")
	// ErrNoSuccess is returned when the recaptcha request was not successful.
	ErrNoSuccess = errors.New("request was not successful")
	// ErrScore is returned when the calculated score is below the required score (V3 only).
	ErrScore = errors.New("request was below the required score")
	// ErrAction is returned when the action is not the required one (V3 only).
	ErrAction = errors.New("wrong action")
)

func verify(secret string, r *http.Request) ([]byte, error) {
	response := r.FormValue("g-recaptcha-response")
	if response == "" {
		return nil, ErrNoCaptcha
	}

	remoteip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	verifyData := url.Values{
		"secret":   {secret},   // private key
		"response": {response}, // response from the client to verify
		"remoteip": {remoteip}, // client ip (optional)
	}

	jsonResp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", verifyData)
	if err != nil {
		return nil, err
	}
	defer jsonResp.Body.Close()

	body, err := ioutil.ReadAll(jsonResp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
