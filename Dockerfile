FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
ENV USER=prom
ENV UID=10001
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"WORKDIR $GOPATH/src/mypackage/myapp/

WORKDIR /app

COPY . .

RUN GOOS=linux go build -ldflags="-w -s" -o /go/bin/paperless_exporter

# Final image
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/paperless_exporter /go/bin/paperless_exporter
USER appuser:appuser
ENTRYPOINT ["/go/bin/paperless_exporter"]
