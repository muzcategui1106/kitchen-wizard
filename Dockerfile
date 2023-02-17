
FROM golang:1.20.0-alpine as build

RUN apk add g++ && apk add make

# Download necessary Go modules
WORKDIR /app

# pre install dependnecies so that changes in code does not trigger a reinstall as the layer is already built
RUN go install github.com/bufbuild/buf/cmd/buf@v1.9.0

# copy the project
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY pkg pkg
COPY Makefile Makefile

# download grpc deps deps
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# generate grpc deps
RUN make grpc-build

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
COPY --from=build /app/bin/kitchen-wizard /bin/kitchen-wizard






