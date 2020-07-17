FROM golang:1.12 as builder
WORKDIR /go/src/github.com/tintinnabulate/registration-webapp
COPY . /go/src/github.com/tintinnabulate/registration-webapp

ENV GO111MODULE=on
# download go modules
RUN go mod download

# build test binary named main.test
RUN go test -c -o main.test ./...

FROM google/cloud-sdk:alpine as alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root
# copy the test binary into the docker container
COPY --from=builder /go/src/github.com/tintinnabulate/registration-webapp/main.test .

WORKDIR /root
CMD ["./main.test -cover -race"]
