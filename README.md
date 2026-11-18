# Hello World API

A small Go web service exposing a single endpoint, `GET /hello-world`, built with the standard library only.

## Endpoint

`GET /hello-world?name=<name>`

| Request                         | Status | Body                           |
|---------------------------------|--------|--------------------------------|
| `/hello-world?name=Alice`       | 200    | `{"message":"Hello Alice"}`    |
| `/hello-world?name=Zane`        | 400    | `{"error":"Invalid Input"}`    |
| `/hello-world` or `?name=`      | 400    | `{"error":"Invalid Input"}`    |

A name is valid if its first letter is A–M or a–m.

## Project structure

```
cmd/server/            entry point: wiring, HTTP server, graceful shutdown
internal/controller/   HTTP layer: reads the query param, maps results to JSON + status codes
internal/service/      business rules: validation and greeting logic (no HTTP)
internal/repository/   data layer: supplies the greeting prefix (in-memory for now)
```

Dependencies point one way: `controller → service → repository`. Each layer depends on an interface, so each one can be tested with a stub.

## Requirements

- Go 1.22 or later (the router uses Go 1.22 method patterns such as `"GET /hello-world"`).

## Running the application

```bash
go run ./cmd/server          # or: make run
```

The server listens on port `8080`. Set `PORT` to change it:

```bash
PORT=9090 go run ./cmd/server
```

Try it:

```bash
curl -i "http://localhost:8080/hello-world?name=alice"
curl -i "http://localhost:8080/hello-world?name=Zane"
curl -i "http://localhost:8080/hello-world"
```

## Running the tests

```bash
go test ./...                # or: make test (runs with -race)
go test -cover ./...         # coverage summary
```

The tests cover:
- **Service**: table-driven tests for valid names, the A–M / N–Z boundaries, and edge cases.
- **Controller**: HTTP tests through the router with `httptest`, checking status codes, `Content-Type` and JSON bodies, plus the 405 and 500 paths.
- **Repository**: the default greeting prefix.

## Assumptions

1. **Case-insensitive.** `alice` and `Alice` are treated the same.
2. **Whitespace is trimmed.** `?name=%20Bob%20` returns `Hello Bob`. A name made only of whitespace counts as empty.
3. **The name must start with an ASCII letter.** Names starting with a digit, a symbol, an accented letter (`Élodie`) or an emoji return `400 Invalid Input`, since the spec only defines the English alphabet. Only the first character is checked, so characters after it can be anything.
4. **The name is echoed back as sent (after trimming)**, with its original casing. JSON encoding escapes it safely.
5. **Repeated parameters**: `?name=Alice&name=Zane` uses the first value.
6. **Only `GET` is supported.** Other methods return `405 Method Not Allowed` from the Go router.
7. **Unexpected errors** return `500` with `{"error":"Internal Server Error"}`. This can't happen with the current in-memory repository, but the controller handles it.
8. The repository is in-memory. It is there to show the layering and to make the greeting source easy to swap (e.g. config or database) without touching the business logic.
