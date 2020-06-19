#Start from the golang base image. we can also use golang:latest
FROM golang:1.14 AS builder

LABEL maintainer="Shohidul bari <shohidulbari18@gmail.com>"
ENV GOPROXY="direct"

WORKDIR /go/src/github.com/sbr35/wallets-users/
COPY . /go/src/github.com/sbr35/wallets-users/
RUN go mod download
EXPOSE 8080
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
#RUN go build -o app .
#CMD ["go", "run", "main.go"]
CMD go get github.com/pilu/fresh && \
    fresh;

#This is for Production to reduce the image size#
#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#COPY --from=builder /go/src/github.com/sbr35/wallets-users/app .
#CMD ["./app"]
