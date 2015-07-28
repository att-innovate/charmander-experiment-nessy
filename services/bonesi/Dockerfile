FROM phusion/baseimage:0.9.16

ENV HOME /root
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -qq

RUN apt-get -y install libpcap-dev libnet1-dev automake1.10 make git

RUN git clone https://github.com/Markus-Go/bonesi \
		&& cd bonesi \
		&& ./configure \
		&& make

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENTRYPOINT ["/bonesi/src/bonesi", "-i", "/bonesi/50k-bots"]


