FROM alpine:3.21 AS cpp-builder

RUN apk add --no-cache alpine-sdk cmake ninja grpc grpc-dev protobuf-dev

COPY rop-algorithm/core/ /build/core/
COPY rop-algorithm/solver/proto/ /build/solver/proto/

WORKDIR /build/core

RUN cmake -S . -B /build/out -G Ninja -DCMAKE_BUILD_TYPE=Release \
 && cmake --build /build/out --target solver --parallel

FROM golang:1.26-alpine AS go-builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app/rop-backend

COPY rop-algorithm/ /app/rop-algorithm/
COPY rop-backend/go.mod rop-backend/go.sum ./
RUN go mod download

COPY rop-backend/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/rop-backend-bin .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata grpc protobuf

RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

COPY --from=go-builder /app/rop-backend-bin ./rop-backend
COPY --from=cpp-builder /build/out/solver ./solver

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 3000

ENTRYPOINT ["/app/rop-backend"]
