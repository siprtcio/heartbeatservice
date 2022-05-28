FROM golang:1.18.2-alpine as builder
WORKDIR /go/src/github.com/siprtcio/heartbeatservice
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /
COPY --from=builder /go/src/github.com/siprtcio/heartbeatservice/app .
COPY --from=builder /go/src/github.com/siprtcio/heartbeatservice/conf /conf

HEALTHCHECK --timeout=5s --interval=3s --retries=3 CMD curl --fail http://localhost:8080/v1/health || exit 1

EXPOSE 8080

CMD ["/app"]
