##
# Build stage
##
FROM golang:1.24.3-alpine AS builder
RUN apk --no-cache add tzdata

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  small-user

WORKDIR /app

# Install dependencies first for caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the application code and build it
COPY main.go ./
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /postlog .

##
# Final stage
##
FROM scratch

ENV PORT=80
EXPOSE 80

# Copy over files for timezone and SSL support
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy over user and group files
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER small-user:small-user

# Copy over the binary from the builder
COPY --from=builder /postlog .

CMD ["./postlog"]