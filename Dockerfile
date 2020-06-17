FROM golang:1.14.1-buster

ENV APP_HOME /app

WORKDIR ${APP_HOME}

COPY . .

CMD ["./go-web-cloner"]

