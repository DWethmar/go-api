FROM golang:1.15

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/dwethmar/go-api
WORKDIR /go/src/github.com/dwethmar/go-api


RUN ./server

EXPOSE 8080
