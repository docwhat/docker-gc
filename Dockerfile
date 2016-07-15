FROM scratch
MAINTAINER Christian Höltje <docwhat@gerf.org>

ENV COLUMNS=80

COPY docker-gc_linux_amd64 /docker-gc

ENTRYPOINT ["/docker-gc"]
