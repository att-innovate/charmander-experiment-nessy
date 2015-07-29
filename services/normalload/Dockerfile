FROM phusion/baseimage:0.9.17


RUN mkdir -p /normalload/queryfiles
COPY ./normalload /normalload/
COPY ./queryfiles/medium_1500 /normalload/queryfiles/

ENTRYPOINT ["/normalload/normalload"]
