##
# Build stage
##
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Install dependencies first for caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the application code and build it
COPY . .
RUN CGO_ENABLED=0 go build -o /postlog .

##
# Final stage
##
FROM scratch

# Note - there's no SSL or timezone support in this image

COPY --from=builder /postlog .

CMD ["./postlog"]