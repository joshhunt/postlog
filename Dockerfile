FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /main .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /main .

CMD ["./main"]