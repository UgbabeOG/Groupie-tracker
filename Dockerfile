FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/groupie-tracker .

FROM alpine:3.19 AS runtime
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/groupie-tracker ./groupie-tracker
COPY --from=builder /src/static ./static
COPY --from=builder /src/templates ./templates
COPY --from=builder /src/data ./data

EXPOSE 8080
CMD ["./groupie-tracker"]
