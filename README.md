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
$ bun run scripts/watch
```

### Generate

Generate files

```sh
$ bun run scripts/generate
```

### Build

Build the Go binary

```sh
$ bun run scripts/build
```

## UNLICENSE

This project is (un)licensed under the [Unlicense](UNLICENSE).
