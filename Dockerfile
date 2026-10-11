
FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /finops-webhook .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /finops-webhook /app/finops-webhook

EXPOSE 8443

ENTRYPOINT ["/app/finops-webhook"]
