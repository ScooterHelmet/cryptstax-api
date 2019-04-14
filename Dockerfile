FROM golang:1.12.4

WORKDIR /bin
RUN mkdir cryptstax-api
COPY . /bin/cryptstax-api

WORKDIR /bin/cryptstax-api
CMD ["./cryptstax-api"]
EXPOSE 8000