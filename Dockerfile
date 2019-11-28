FROM golang:alpine AS builder

RUN apk add --no-cache \
  git

COPY main.go .

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-w -s" -o /rallytics

# --8<-------------
FROM scratch

COPY --from=builder /rallytics /rallytics

ENTRYPOINT ["/rallytics"]
