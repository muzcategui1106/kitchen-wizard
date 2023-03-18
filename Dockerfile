
FROM golang:1.20.0-alpine as build

RUN apk add g++ && apk add make

# Download necessary Go modules
WORKDIR /app

# copy the project
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY pkg pkg
COPY scripts scripts
COPY Makefile Makefile

# setting path
RUN export PATH="$PATH:$GOPATH/bin"

# get dependencies
RUN go mod download

# build code
RUN make go-build

#################################

# final image\
# TODO improve this we want to use scratch but we cant because it complains about not finding the executable
# I believe we just gonna enable CGO during the build
FROM golang:1.20.0-alpine AS final

WORKDIR /
COPY --from=build  /app/swagger-ui /swagger-ui
COPY --from=build /app/bin/api /api






