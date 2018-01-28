FROM alpine:3.7

EXPOSE 8080

ADD ./api2html /etc/api2html/api2html

WORKDIR /etc/api2html/

ENTRYPOINT [ "/etc/api2html/./api2html" ]

CMD [ "-h" ]
