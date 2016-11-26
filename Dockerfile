FROM alpine:latest

RUN apk --update add ca-certificates
EXPOSE 8080

ENTRYPOINT ["./mws-bridge"]

COPY _build/mws-bridge .
