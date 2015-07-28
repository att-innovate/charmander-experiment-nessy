FROM phusion/baseimage:0.9.16


COPY ./normalload /usr/bin/
RUN mkdir /data
COPY ./files/tiny_20 /data/

ENTRYPOINT ["/usr/bin/normalload"]
