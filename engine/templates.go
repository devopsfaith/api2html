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

    <h2>Request context parameters (<tt>Context.params</tt>)</h2>
    <div class="response">
        {{ #Context.params }}
            <pre>{{ . }}</pre>
        {{ /Context.params }}
        {{ ^Context.params }}
            <p>No context parameters set.</p>
        {{ /Context.params }}
    </div>
    <h2>Request context keys (<tt>Context.keys</tt>)</h2>
    <div class="response">

        {{ #Context.keys }}
            <pre>{{ . }}</pre>
        {{ /Context.keys }}
        {{ ^Context.keys }}
            <p>No context keys set.</p>
        {{ /Context.keys }}

    </div>
    <h2>Request parameters (<tt>Params</tt>)</h2>
    <div class="response">
        {{ #Params }}
        <pre>{{ . }}</pre>
        {{ /Params }}
         {{ ^Params }}
            <p>This page didn't set any parameters in the URL.</p>
        {{ /Params }}
    </div>
    <h2>Extra data from config(<tt>Extra</tt>)</h2>
    <div class="response">
        {{ #Extra }}
        <pre>{{ . }}</pre>
        {{ /Extra }}
        {{ ^Extra }}
            <p>The configuration file does not add any extra data.</p>
        {{ /Extra }}
    </div>
    <h2>Backend response</h2>
    <div class="response">
        <h3>Response when object (<tt>Data</tt>)</h3>
        {{ #Data }}
        <pre>{{ . }}</pre>
        {{ /Data }}
        {{ ^Data }}
            <p>The backend response did not return an object.</p>
        {{ /Data }}

        <h3>Response when array (<tt>Array</tt>)</h3>
        {{ #Array }}
        <pre>{{ . }}</pre>
        {{ /Array }}
        {{ ^Array }}
            <p>The backend response did not return an array or configuration does not set <tt>isArray</tt>.</p>
        {{ /Array }}
    </div>
</div>
<style type="text/css">
    .api2html-debug {
        background-color: #f1f1f1;
        border: 1px solid #666;
        color: #333;
        margin:2rem;
    }
    .api2html-debug .response {
        padding: 1em;
    }
    .api2html-debug pre, .api2html-debug strong {
        color: #cb2027;
        font-family: monospace;
    }

    .api2html-debug h1 {
        text-align: center;
    }

    .api2html-debug h1, .api2html-debug h2, .api2html-debug h3 {
        margin: 0;
        background-color: #e0e0e0;
        color: #cb2027;
        padding:0.5em;
    }
</style>`
)
