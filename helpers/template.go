package helpers

import (
	"bytes"
	"html/template"
)

// RenderRedirectTemplate renders the redirect template
func RenderRedirectTemplate(method, url string, data map[string]string) (string, error) {
	tmp, _ := template.New("html").Parse(redirectTemplate)
	var tpl bytes.Buffer
	_ = tmp.Execute(&tpl, map[string]interface{}{
		"method": method,
		"payURL": url,
		"data":   data,
	})
	return tpl.String(), nil
}

const redirectTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Forwarding to secure payment provider</title>
    <style>
        .text-center {
            text-align: center;
        }

        .mt-2 {
            margin-top: 2em;
        }

        .spinner {
            margin: 100px auto 0;
            width: 70px;
            text-align: center;
        }

        .spinner > div {
            width: 18px;
            height: 18px;
            background-color: #333;
            border-radius: 100%;
            display: inline-block;
            -webkit-animation: sk-bouncedelay 1.4s infinite ease-in-out both;
            animation: sk-bouncedelay 1.4s infinite ease-in-out both;
        }

        .spinner .bounce1 {
            -webkit-animation-delay: -0.32s;
            animation-delay: -0.32s;
        }

        .spinner .bounce2 {
            -webkit-animation-delay: -0.16s;
            animation-delay: -0.16s;
        }

        @-webkit-keyframes sk-bouncedelay {
            0%, 80%, 100% { -webkit-transform: scale(0) }
            40% { -webkit-transform: scale(1.0) }
        }

        @keyframes sk-bouncedelay {
            0%, 80%, 100% {
                -webkit-transform: scale(0);
                transform: scale(0);
            } 40% {
                  -webkit-transform: scale(1.0);
                  transform: scale(1.0);
              }
        }
    </style>
</head>
<body>
<div class="spinner">
    <div class="bounce1"></div>
    <div class="bounce2"></div>
    <div class="bounce3"></div>
</div>
<form class="text-center mt-2" method="{{ .method }}" action="{{ .payURL }}">
    <p>Forwarding to secure payment provider.</p>
    <p>
        If you are not automatically redirected to the payment website with in
        <span id="countdown">10</span>
        seconds...
    </p>
    {{ range $key, $value := .data }}
        <input type="hidden" name="{{ $key }}" value="{{ $value }}">
    {{ end }}

    <button type="submit">Click here</button>
</form>
<script>
    // Total seconds to wait
    let seconds = 10;

    function submitForm() {
        document.forms[0].submit();
    }

    function countdown() {
        seconds = seconds - 1;
        if (seconds <= 0) {
            // submit the form
            submitForm();
        } else {
            // Update remaining seconds
            document.getElementById("countdown").innerHTML = seconds;
            // Count down using javascript
            window.setTimeout("countdown()", 1000);
        }
    }

    // Run countdown function
    countdown();
</script>
</body>
</html>
`
