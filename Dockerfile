FROM ubuntu:trusty

MAINTAINER Stefan Arentz <stefan@arentz.ca>

ADD moz-csp-collector /usr/local/bin/moz-csp-collector
ADD moz-csp-collector.sh /usr/local/bin/moz-csp-collector.sh

EXPOSE 8080

CMD ["/usr/local/bin/moz-csp-collector.sh"]
