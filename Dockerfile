FROM golang

LABEL author="emrez"
LABEL version="1.1"

RUN mkdir /app

WORKDIR /app

COPY ./app ./app
COPY ./domain ./domain
COPY ./controller ./controller
COPY ./infra ./infra
COPY ./cache ./cache
COPY ./server ./server
COPY ./.config ./.config
COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /uygulama

CMD ["/uygulama"]

EXPOSE 8080
