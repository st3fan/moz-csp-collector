FROM ubuntu:trusty

MAINTAINER Stefan Arentz <stefan@arentz.ca>

ADD moz-csp-collector /usr/local/bin/moz-csp-collector

EXPOSE 8080

CMD ["/usr/local/bin/moz-csp-collector"]
