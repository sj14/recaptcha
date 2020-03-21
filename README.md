# reCAPTCHA for Go

## Installation

```bash
go get -u github.com/sj14/recaptcha
```

## Usage

The result will always return an error when any kind of failure occured, thus only checking the error is sufficient.

### reCAPTCHA V2

```go
result, err := recaptcha.VerifyV2(recaptchaSecret, httpRequest)
if err != nil {
    log.Fatalf("recaptcha failed: %v\n", err)
}
// result is not necessary to check, only required for more details
```

### reCAPTCHA V3

Without options:

```go
result, err := recaptcha.VerifyV3(recaptchaSecret, httpRequest)
if err != nil {
    log.Fatalf("recaptcha failed: %v\n", err)
}
// result is not necessary to check, only required for more details
```

With options:

```go
// The default action is empty ("") and thus not checked.
// The default minimum required score is 0.5
result, err := recaptcha.VerifyV3(recaptchaSecret, httpRequest, recaptcha.Action("register"), recaptcha.MinScore(0.7))
if err != nil {
    log.Fatalf("recaptcha failed: %v\n", err)
}
// result is not necessary to check, only required for more details
```
