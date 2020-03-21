# reCAPTCHA for Go

## Installation

```bash
go get -u github.com/sj14/recaptcha
```

## Usage

The result will always return an error when any kind of failure occured, thus only checking the error is sufficient.

### reCAPTCHA V2

#### HTML

See https://developers.google.com/recaptcha/docs/display#automatically_render_the_recaptcha_widget

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

### reCAPTCHA V3

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
