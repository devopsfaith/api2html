# Example blog

The following code shows a very basic example of API2HTML serving pages. In order to run it:

	cd examples/blog
	api2html serve -d -c config.json

This will start the server and you will be able to navigate the following pages:

	- [Home](http://localhost:8080/)
	- [Post](http://localhost:8080/posts/1)
	- [robots.txt](http://localhost:8080/robots.txt)
	- [sitemap.xml](http://localhost:8080/sitemap.xml)
	- [hello.txt](http://localhost:8080/hello.txt): An example of `text/plain` content
	- [404 page](http://localhost:8080/idontexist)
	- 500 page: You need to break the API response to see it

See the `config.json` to understand how this works and the `tmpl` folder to see how `{{data}}` is injected.

The template system is [Mustache](https://mustache.github.io)
