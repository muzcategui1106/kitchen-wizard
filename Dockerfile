FROM golang:1.20.0-alpine as build

RUN apk add g++ && apk add make
RUN apk update && apk add bash
RUN apk update && apk add rsync

# setting path
RUN export PATH="$PATH:$GOPATH/bin"

# Download necessary Go modules
WORKDIR /app

# copy scripts so we can just warm up swag dependencies so that minor code changes do not retrigger a full build all the time
COPY scripts scripts
RUN scripts/get-swag-dependencies.sh

# copy go.mod and go.sum so we can just warm up swag dependencies so that minor code changes do not retrigger a full build all the time
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# copy the project code
COPY main.go main.go
COPY pkg pkg
COPY scripts scripts
COPY Makefile Makefile

# build code
RUN make go-build

#################################

# final image\
# TODO improve this we want to use scratch but we cant because it complains about not finding the executable
# I believe we just gonna enable CGO during the build
FROM golang:1.20.0-alpine AS final

WORKDIR /
COPY --from=build /app/bin/api /api






