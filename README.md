# Go Minimal Starter

Minimal but opinionated starter project for Go applications

## Tech Stack

- **Go** - Main programming language used for backend development
  - **wgo** - Live reload for Go apps
  - **koanf** - Simple configuration library for Go
- **htmx** - Access modern browser features directly from HTML
  - **templ** - Template engine for Server-Side Rendering
- **Tailwind CSS** - Utility-first CSS framework
  - **daisyUI** - Component library for Tailwind CSS
- **Bun** - High-performance JavaScript runtime and build tool
- **PostgreSQL** - Open-source relational database
  - **sqlc** - Generate type-safe Go code from SQL
  - **pgx** - PostgreSQL driver and toolkit for Go
  - **dockertest** - Run Docker containers for integration tests

## Getting Started

### Prerequisite

- `go` : `^1.24`
- `bun` : `^1.2`

### Init

Install dependencies

```sh
$ go mod download && bun i
```

### Watch

Live reloading

```sh
$ bun run watch
```

### Generate

Generate files

```sh
$ bun run generate
```

### Build

Build the Go binary

```sh
$ bun run build
```

### Test

Run tests

```sh
$ bun run test
```

## License

This project is released into the public domain under the [Unlicense](UNLICENSE).
