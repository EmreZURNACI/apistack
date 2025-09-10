FROM golang:tip-bullseye

LABEL author="emrez"
LABEL version="1.0"

RUN mkdir /app

WORKDIR /app

COPY ./app ./app
COPY ./domain ./domain
COPY ./controller ./controller
COPY ./infra ./infra
COPY ./server ./server
COPY ./.env ./.env
COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum


RUN go mod tidy
RUN go mod download

CMD ["go","run","main.go"]

EXPOSE 8080:8080
