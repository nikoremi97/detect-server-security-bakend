FROM alpine
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
WORKDIR /app
COPY main /app/main
ENTRYPOINT /app/main