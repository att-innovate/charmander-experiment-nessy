FROM phusion/baseimage:0.9.16

ADD ./load_random /usr/bin/load_random
RUN chmod +x /usr/bin/load_random

ENTRYPOINT ["/usr/bin/load_random"]
