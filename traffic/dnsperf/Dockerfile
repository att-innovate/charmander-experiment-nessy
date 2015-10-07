FROM phusion/baseimage:0.9.17

COPY ./queryfiles/tiny_20 /queryfiles/ 
COPY ./queryfiles/small_100 /queryfiles/ 
COPY ./queryfiles/medium_1500 /queryfiles/ 


RUN apt-get -y update && apt-get install -y dnsutils build-essential wget bind9utils libbind-dev libkrb5-dev libssl-dev libcap-dev libxml2-dev


RUN wget ftp://ftp.nominum.com/pub/nominum/dnsperf/2.0.0.0/dnsperf-src-2.0.0.0-1.tar.gz \
		&& tar xfvz dnsperf-src-2.0.0.0-1.tar.gz \
		&& cd dnsperf-src-2.0.0.0-1 \
		&& ./configure \
		&& make \
		&& make install 

ENTRYPOINT ["/usr/local/bin/dnsperf"]
