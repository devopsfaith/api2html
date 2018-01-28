api2html
====

Template Render as a Service

## Dependencies

Dependencies are managed by `dep`. Install it with:

	$ make prepare

## Build

	$ make

## Run

	$ ./api2html -h
	Template Render As A Service

	Usage:
	  api2html [command]

	Available Commands:
	  generate    Generate the final api2html templates.
	  serve       Run the api2html server.

	Use "api2html [command] --help" for more information about a command.

## Run the engine

	$ ./api2html run -h
	Run the api2html server.

	Usage:
	  api2html serve [flags]

	Aliases:
	  serve, run, server, start


	Examples:
	api2html serve -d -c config.json

	Flags:
	  -c, --config string   Path to the configuration filename (default "config.json")
	  -d, --devel           Enable the devel

## Generator

	$ ./api2html generate -h
	Generate the final api2html templates.

	Usage:
	  api2html generate [flags]

	Aliases:
	  generate, create, new


	Examples:
	api2html generate -d -c config.json

	Flags:
	  -i, --iso string    (comma-separated) iso code of the site to create (default "*")
	  -p, --path string   Base path for the generation (default ".")
	  -r, --reg string    regex filtering the sources to move to the output folder (default "ignore")

## Hot template reload

	$ curl -X PUT -F "file=@/path/to/tmpl.mustache" -H "Content-Type: multipart/form-data" \
	http://localhost:8080/template/<TEMPLATE_NAME>

## Docker build

	$ make docker

## Docker run

	$ docker run -it --rm -p8080:8080 -v $PWD/config.json:/etc/api2html/config.json api2html -d -c /etc/api2html/config.json