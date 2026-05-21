FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app/rop-backend

COPY rop-algorithm/ /app/rop-algorithm/
COPY rop-backend/go.mod rop-backend/go.sum ./
RUN go mod download

COPY rop-backend/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/rop-backend-bin .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

COPY --from=builder /app/rop-backend-bin ./rop-backend
COPY rop-backend/solver-bin/solver ./solver

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 3000

ENTRYPOINT ["/app/rop-backend"]
