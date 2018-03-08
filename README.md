
![api2html](https://raw.githubusercontent.com/devopsfaith/api2html.com/master/img/API2HTML-on-white.png)

[![Build Status](https://travis-ci.org/devopsfaith/api2html.svg?branch=master)](https://travis-ci.org/devopsfaith/api2html) [![Go Report Card](https://goreportcard.com/badge/github.com/devopsfaith/api2html)](https://goreportcard.com/report/github.com/devopsfaith/api2html) [![Coverage Status](https://coveralls.io/repos/github/devopsfaith/api2html/badge.svg?branch=master)](https://coveralls.io/github/devopsfaith/api2html?branch=master) [![GoDoc](https://godoc.org/github.com/devopsfaith/api2html?status.svg)](https://godoc.org/github.com/devopsfaith/api2html)

### On the fly HTML generator from API data

API2HTML is a web server that renders [Mustache](http://mustache.github.io/) templates and injects them your API data. This allows you to build websites by just declaring the API sources and writing the template view.

## How does it work?
To create pages that feed from a backend you just need to add in the configuration file the URL patterns the server will listen to. Let's imagine we want to offer URLs like `/products/13-inches-laptops` where the second part is a variable that will be sent to the API:

    ...
    "pages":[
    {
        "name": "products",
        "URLPattern": "/products/:category",
        "BackendURLPattern": "http://api.company.com/products/:category",
        "Template": "products_list",
        "CacheTTL": "3600s",
        "extra": {
            "promo":"Black Friday"
        }
    },
    ...

The `Template` setting will look for the file `tmpl/products_list.mustache` and the response of the BackendURLPattern call will be injected in the variable `data`. An example of how you could use it:

    <h1>Products for sale</h1>
    <p>Take advantage of the {{extra.promo}}!</p>
    <table>
        {{#data}}
            <tr>
                <td>{{name}}</td>
                <td>{{price}}</td>
            </tr>
        {{/data}}

        {{^data}}
           <tr>
                <td colspan="2">There are no products in this category</td>
            </tr>
        {{/data}}
    </table>

You probably guessed it already, but in this scenario the backend would be returning a response like this:

    // http://api.company.com/products/13-inches-laptops
    {
        [
            { "name": "13-inch MacBook Air", "price": "$999.00" },
            { "name": "Lenovo ThinkPad 13", "price": "$752.00" },
            { "name": "Dell XPS13", "price": "$925.00" }
        ]
    }


## Install

When you install `api2html` for the first time you need to download the dependencies, automatically managed by `dep`. Install it with:

    $ make prepare

Once all dependencies are installed just run:

    $ make

## Run
Once you have successfully compiled API2HTML in your platform the binary `api2html` will exist in the folder. Execute it as follows:

    $ ./api2html -h
    Template Render As A Service

    Usage:
      api2html [command]

    Available Commands:
      generate    Generate the final api2html templates.
      serve       Run the api2html server.

    Use "api2html [command] --help" for more information about a command.

### Run the engine

    $ ./api2html run -h
    Run the api2html server.

    Usage:
      api2html serve [flags]

    Aliases:
      serve, run, server, start

    Examples:
    api2html serve -d -c config.json -p 8080

    Flags:
      -c, --config string   Path to the configuration filename (default "config.json")
      -d, --devel           Enable the devel
      -p, --port int        Listen port (default 8080)

### Generator
The generator allows you to create multiple mustache files using templating. That's right create templates with templates!

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

### Hot template reload

    $ curl -X PUT -F "file=@/path/to/tmpl.mustache" -H "Content-Type: multipart/form-data" \
    http://localhost:8080/template/<TEMPLATE_NAME>

## Building and running with Docker
To build the project with Docker:

    $ make docker

And run it as follows:

    $ docker run -it --rm -p8080:8080 -v $PWD/config.json:/etc/api2html/config.json api2html -d -c /etc/api2html/config.json
