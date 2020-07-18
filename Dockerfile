FROM golang:1.12 as builder
WORKDIR /app
COPY go.mod go.sum ./

ENV GO111MODULE=on
# download go modules
RUN go mod download

# copy source from current directory to working directory inside the container
COPY . .

# build test binary named main.test
RUN CGO_ENABLED=1 GOOS=linux GOPROXY=https://proxy.golang.org go test -c -o main.test -cover ./...

############ start a new stage from scratch ############

FROM google/cloud-sdk:alpine as alpine
RUN apk --no-cache add ca-certificates
# gcompat is necessary for libc
RUN apk update && apk add --no-cache gcompat
WORKDIR /root/

# copy the test binary into the docker container
COPY --from=builder /app/main.test .

# copy necessary environment files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/locales ./locales
COPY --from=builder /app/config.example.json .
# TODO: can we remove the need for this?
COPY --from=builder /app/fanjoula.json .

# run testsuite
CMD ["./main.test -test.v"]
