package engine

var (
	default404Tmpl = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<title>Page not found</title>
</head>
<body class="text-center">
	<h1 class="my-5">Page not found!</h1>
	<p>The page you are looking for is not hosted in this site</p>
	<p>You might want to customize this file by editing <code>static/404</code></p>
</body>`

	default500Tmpl = `<!DOCTYPE html>
<html lang="es">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<title>Bummer!</title>
</head>
<body class="text-center">
	<h1 class="my-5">API response went wild!</h1>
	<p>The response from the API was weird an unable to process it.</p>
	<p>You might want to customize this file by editing <code>static/500</code></p>
</body>`

	debuggerTmpl = `<div class="api2html-debug">
    <h1>API2HTML Debugger</h1>
    <p class="response">Page generated at <strong>{{ Helper.Now }}</strong></p>
    <h2>Response context</h2>
    <div class="response">
    <pre>{{ String }}</pre>
    </div>

    <h2>Request context params</h2>
    <div class="response">
        {{ #Context.params }}
            <pre>{{ . }}</pre>
        {{ /Context.params }}
        {{ ^Context.params }}
            <p>No context parameters set</p>
        {{ /Context.params }}
    </div>
    <h2>Request context keys</h2>
    <div class="response">

        {{ #Context.keys }}
            <pre>{{ . }}</pre>
        {{ /Context.keys }}
        {{ ^Context.keys }}
            <p>No context keys set</p>
        {{ /Context.keys }}

    </div>
    <h2>Request params</h2>
    <div class="response">
        {{ #Params }}
        <pre>{{ . }}</pre>
        {{ /Params }}
    </div>
    <h2>Extra data</h2>
    <div class="response">
        {{ #Extra }}
        <pre>{{ . }}</pre>
        {{ /Extra }}
    </div>
    <h2>Backend response</h2>
    <div class="response">
        <h3>As object</h3>
        {{ #Data }}
        <pre>{{ . }}</pre>
        {{ /Data }}

        <h3>As array</h3>
        {{ #Array }}
        <pre>{{ . }}</pre>
        {{ /Array }}
    </div>
</div>
<style type="text/css">
    .api2html-debug {
        background-color: #1c1c1d;
        border: 1px solid #000;
        color: #fff;
        font-family: Arial, sans-serif;
    }
    .api2html-debug .response {

        padding: 1em;
    }
    .api2html-debug pre, .api2html-debug strong {
        color: #00FF00;
        font-family: monospace;
    }

    .api2html-debug h1, h2, h3 {
        margin: 0;
        background-color: #333;
        color: #fff;
        padding:0.5em;
        border-left: 4px solid #FF473B;
    }
</style>`
)
