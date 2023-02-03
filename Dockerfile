FROM golang:1.20.0-alpine as build

RUN apk add g++ && apk add make

# Download necessary Go modules
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY pkg pkg
COPY Makefile Makefile

# download deps
RUN go mod download

# build code
RUN make go-build


# final image
FROM scratch AS final

COPY --from=build /app/bin/kitchen-wizard /bin/kitchen-wizard
ENTRYPOINT [ "bin/kitchen-wizard" ]






