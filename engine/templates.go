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

	default405Tmpl = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<title>Method not supported</title>
</head>
<body class="text-center">
	<h1 class="my-5">Method not allowed!</h1>
	<p>The requested HTTP method is not supported</p>
	<p>You might want to customize this file by editing <code>static/405</code></p>
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
)
