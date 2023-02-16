FROM golang:1.20.0-alpine

WORKDIR /app

# download grpc deps deps
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN apk add --no-cache protobuf git \
  && go install github.com/golang/protobuf/protoc-gen-go@latest \
  && cp /go/bin/protoc-gen-go /usr/bin/
RUN export PATH="$PATH:$(go env GOPATH)/bin"

COPY go.mod go.mod

RUN go get github.com/googleapis/googleapis@v0.0.0-20230216195746-9de3a8da8b84
RUN go get github.com/grpc-ecosystem/grpc-gateway@v1.16.0
RUN git clone https://github.com/protocolbuffers/protobuf --depth=1
RUN mv protobuf $GOPATH/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20230216195746-9de3a8da8b84/google/
ENV INCLUDE_GOOGLE_APIS="$GOPATH/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20230216195746-9de3a8da8b84"
ENV INCLUDE_GRPC_GATEWAY="$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0"

COPY scripts scripts
RUN chmod +x scripts/build-grpc.sh
# ENTRYPOINT [ /app/scripts/build-grpc.sh ]

