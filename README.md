# Postlog

A lightweight HTTP service that logs incoming payloads to the console. Send JSON data via POST requests or query parameters via GET requests, and the service will log all received data with structured output.

Originally created to receive events from Apple Shortcuts and save them to Grafana Loki. Even this is probably
overkill for my original needs, but it was a fun little exercise.

## Quick Start

### Development

```bash
go run main.go
```

The service runs on port 8080 by default.

### Docker

Run the latest image:

```bash
docker run -p 8080:80 ghcr.io/joshhunt/postlog:latest
```

Or build and run locally:

```bash
docker build -t postlog .
docker run -p 8080:80 postlog
```

## API Endpoints

- `GET /[path]` - Log query parameters from any path
- `POST /[path]` - Log JSON payload from any path (Content-Type: application/json)
- `GET /health` - Health check endpoint

All requests to any path will be logged with the path included as `payload_name` and the HTTP method as `payload_method`.

### Examples

#### GET Request

```bash
$ curl "http://localhost:8080/annotation/ac-on?room=bedroom"
```

Response:

```json
{
  "room": "bedroom",
  "payload_method": "GET",
  "payload_name": "annotation/ac-on"
}
```

Server console:

```json
{
  "level": "info",
  "ts": 1748729270.4279022,
  "msg": "received payload",
  "room": "bedroom",
  "payload_name": "annotation/ac-on",
  "payload_method": "GET"
}
```

#### POST Request

```bash
curl -X POST http://localhost:8080/annotation/ac-on \
       -H "Content-Type: application/json" \
       -d '{ "room": "living-room" }'
```

Response:

```json
{
  "payload_method": "POST",
  "payload_name": "annotation/ac-on",
  "room": "living-room"
}
```

Server console:

```json
{
  "level": "info",
  "ts": 1748729640.583546,
  "caller": "handlers/payload_handler.go:85",
  "msg": "received payload",
  "room": "living-room",
  "payload_name": "annotation/ac-on",
  "payload_method": "POST"
}
```
