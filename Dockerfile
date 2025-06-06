FROM golang:1.23 AS builder

RUN mkdir /logs
RUN mkdir /maps
WORKDIR /app

RUN mkdir ./escape-engine
# RUN mkdir ./escape-models
RUN mkdir ./escape-api

COPY ./escape-engine/ ./escape-engine
# COPY ./escape-models/ ./escape-models
COPY ./escape-api/ ./escape-api

RUN CGO_ENABLED=0 GOOS=linux go build -C ./escape-api -o escapeserver

FROM busybox
COPY --from=builder /app/escape-api/escapeserver /app/
COPY ./escape-api/assets /escape-api/assets
COPY --from=builder /logs /logs
COPY --from=builder /maps /maps
EXPOSE 80

ENTRYPOINT ["/app/escapeserver"]
