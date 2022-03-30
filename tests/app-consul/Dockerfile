ARG GO_IMG

FROM $GO_IMG
ARG CWD
ARG GRPCURL_VERSION=v1.8.6
ENV GOFLAGS=-mod=mod
ENV LOG_LEVEL=debug

RUN apk add --update git ca-certificates && \
    cd /go/src && go mod init example.com/module && \
    go get github.com/fullstorydev/grpcurl@${GRPCURL_VERSION} && \
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl

WORKDIR $CWD
COPY . .
RUN go build -v -o /app $CWD/tests/app-consul

CMD ["/app"]

HEALTHCHECK --interval=16s --timeout=2s \
    CMD grpcurl -plaintext -d {} 127.0.0.1:9000 grpc.health.v1.Health.Check | grep -w "SERVING"
