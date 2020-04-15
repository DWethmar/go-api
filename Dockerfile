FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/dwethmar/go-api
WORKDIR /go/src/github.com/dwethmar/go-api

RUN go get ./cmd/api/
RUN go build ./cmd/api/

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get -u github.com/c9s/gomon && \
	gomon -b ./cmd/api/; \
	fi

EXPOSE 8080
