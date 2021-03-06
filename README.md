# reCAPTCHA for Go

[![GoDoc](https://godoc.org/github.com/sj14/recaptcha?status.png)](https://pkg.go.dev/github.com/sj14/recaptcha?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/sj14/recaptcha)](https://goreportcard.com/report/github.com/sj14/recaptcha)
![Action](https://github.com/sj14/recaptcha/workflows/Go/badge.svg)

## Installation

```bash
go get -u github.com/sj14/recaptcha
```

## Usage

The result will always return an error when any kind of failure occured, thus only checking the error is sufficient.

### reCAPTCHA v2

See [reCAPTCHA v2 docs](https://developers.google.com/recaptcha/docs/display)

#### HTML

```html
<html>
  <head>
    <script src="https://www.google.com/recaptcha/api.js" async defer></script>
  </head>
  <body>
    <form action="?" method="POST">
      <div class="g-recaptcha" data-sitekey="<YOUR_SITE_KEY>"></div>
      <br/>
      <input type="submit" value="Submit">
    </form>
  </body>
</html>
```

#### Go

```go
result, err := recaptcha.VerifyV2(recaptchaSecret, httpRequest)
if err != nil {
    log.Fatalf("recaptcha failed: %v\n", err)
}
// result is not necessary to check, only required for more details
```

### reCAPTCHA v3

See [reCAPTCHA v3 docs](https://developers.google.com/recaptcha/docs/v3)

#### HTML

```html
<html>
  <head>
    <script src="https://www.google.com/recaptcha/api.js?render=<YOUR_SITE_KEY>"></script>
    <script>
    grecaptcha.ready(function () {
        grecaptcha.execute('<YOUR_SITE_KEY>', { action: 'register' }).then(function (token) {
            var recaptchaResponse = document.getElementById('g-recaptcha-response');
            recaptchaResponse.value = token;
        });
    });
</script>
</head>
<body>
    <form method="POST" action="/register">
        <input type="hidden" name="g-recaptcha-response" id="g-recaptcha-response">
        <input type="submit" value="Submit">
    </form>
</body>

```

#### Go

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

## Shoulders

This package is based on the work of [Adam Ng](https://www.socketloop.com/tutorials/golang-recaptcha-example) and I highly support his appeal:

> IF you gain some knowledge or the information here solved your programming problem. Please consider donating to the less fortunate or some charities that you like. Apart from donation, planting trees, volunteering or reducing your carbon footprint will be great too.
