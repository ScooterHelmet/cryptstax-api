FROM golang:1.12.4

WORKDIR /bin
RUN mkdir cryptstax-api
COPY . /bin/cryptstax-api

WORKDIR /bin/cryptstax-api/db
RUN cockroach start --insecure --host 127.0.0.1 \
    && cat schema.sql | cockroach sql --insecure
WORKDIR /bin/cryptstax-api
CMD ["./cryptstax-api"]
EXPOSE 8000