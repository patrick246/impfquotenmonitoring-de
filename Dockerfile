ARG GOVERSION=1.15
FROM golang:$GOVERSION as builder
WORKDIR /usr/share/app
ADD . .
RUN go build -o bin --ldflags '-extldflags "-static"' -tags osusergo,netgo ./cmd/...

FROM gcr.io/distroless/static
WORKDIR /
COPY --from=builder /usr/share/app/bin/* /app/

ENTRYPOINT ["/app/server"]
