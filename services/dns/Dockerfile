FROM phusion/baseimage:0.9.17

ENV HOME /root
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -qq

RUN apt-get -y install bind9

COPY files/named.conf.options /etc/bind/named.conf.options
COPY files/named.conf.local /etc/bind/named.conf.local
COPY files/db.homer.ave /etc/bind/zones/db.homer.ave
COPY files/db.192 /etc/bind/zones/db.192

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENTRYPOINT service bind9 start && tail -f /var/log/syslog


