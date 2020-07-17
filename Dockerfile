FROM golang:1.12 as builder
WORKDIR /app
COPY go.mod go.sum ./

ENV GO111MODULE=on
# download go modules
RUN go mod download

# copy source from current directory to working directory inside the container
COPY . .

# build test binary named main.test
RUN go test -c -o main.test ./...

# start a new stage from scratch
FROM google/cloud-sdk:alpine as alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# copy the test binary into the docker container
COPY --from=builder /app/main.test .

# run testsuite
CMD ["./main.test -cover -race"]
