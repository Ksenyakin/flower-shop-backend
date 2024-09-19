FROM ubuntu:latest
LABEL authors="DanKs"

ENTRYPOINT ["top", "-b"]