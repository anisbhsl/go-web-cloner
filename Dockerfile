FROM golang:1.14.1-buster

ENV APP_HOME /app

WORKDIR ${APP_HOME}

COPY . .

RUN go build -o go-web-cloner cmd/go-web-cloner/main.go

ENTRYPOINT ["./go-web-cloner"]

