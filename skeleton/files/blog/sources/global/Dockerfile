FROM devopsfaith/api2html

ADD tmpl/* /etc/api2hmtl/
ADD static/* /etc/api2html/static/
ADD config.json /etc/api2html/config.json

CMD [ "-d", "-c", "/etc/api2hmtl/config.json" ]
