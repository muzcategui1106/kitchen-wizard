FROM golang:1.20.0-alpine

RUN apk add g++ && apk add make

# Download necessary Go modules
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY Makefile Makefile

# download deps
RUN go mod download

# build code
RUN make build-local

ENTRYPOINT [ "bin/kitchen-wizard" ]






