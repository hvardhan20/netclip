# syntax=docker/dockerfile:1

FROM golang:alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
COPY *.html ./
RUN go build -o /go-netclip
EXPOSE 8081
CMD [ "/go-netclip" ]
