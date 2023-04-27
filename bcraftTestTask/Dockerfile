FROM golang:1.19-alpine as builder
WORKDIR /build
COPY . .
COPY /config /config
RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go
FROM alpine:3
COPY --from=builder /config /config
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]
EXPOSE 8000