FROM golang:1.12 as builder
WORKDIR /go/src/github.com/tintinnabulate/registration-webapp
COPY . /go/src/github.com/tintinnabulate/registration-webapp

ENV GO111MODULE=on
# download go modules
RUN go mod download

# build test binary named main.test
RUN go test -c -o main.test ./...

FROM alpine:latest as alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root
# copy the test binary into the docker container
COPY --from=builder /go/src/github.com/tintinnabulate/registration-webapp/main.test .

WORKDIR /tmp
RUN apk add curl
# download google cloud sdk
RUN curl -o /tmp/sdk.tar.gz "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-258.0.0-linux-x86_64.tar.gz" 
RUN tar xzf /tmp/sdk.tar.gz -C /tmp/sdk
RUN export PATH="$PATH:/tmp/sdk/bin"

WORKDIR /root
CMD ["./main.test -cover -race"]
