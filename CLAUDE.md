# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**cles** is a CLI client for Elasticsearch written in Go. It wraps the Elasticsearch API to provide user-friendly commands for index management, cat API, search templates, bulk operations, and reindexing.

## Build & Development Commands

```bash
make build          # Build binary (outputs ./cles)
make test           # Build + run all tests with coverage
make clean          # Remove build artifacts
make install        # Install to $GOPATH/bin
go test -v -run TestFunctionName ./...  # Run a single test
```

CI runs `make test` on every push/PR (see `.github/workflows/test.yaml`).

## Architecture

All source files are in the root package `main` — there are no sub-packages.

- **app.go** — App initialization (`initApp`), version constant, global flags (`--profile`, `--debug`). `Before` hook sets up debug logging and ES client via context.
- **commands.go** — Command tree definition using `urfave/cli/v2`. Each `load*Command()` function returns a `*cli.Command` with subcommands.
- **config.go** — Config loading: TOML file (`$HOME/.config/cles/config.toml`) with env var overrides (`ES_ADDRESS`, `ES_USERNAME`, `ES_PASSWORD`, `ES_SNIFF`). Also creates the `olivere/elastic` client.
- **cat.go** — `cat aliases` / `cat indices` commands + pretty-print formatters with interfaces (`catAliasesSvc`, `catIndicesSvc`) for testability.
- **indices.go** — Index CRUD: create, delete, alias, mapping.
- **search_template.go** — Search template CRUD + render/search.
- **bulk.go** — Bulk index from ndjson files.
- **reindex.go** — Reindex API wrapper.
- **debug.go** — Debug logging setup (writes to stderr when `--debug` is set).
- **main.go** — Entry point; calls `initApp()` and `run()`.

## Key Patterns

- The ES client is passed through `cli.Context` via `context.WithValue` (key: `"client"`). Commands retrieve it with `c.Context.Value("client").(*elastic.Client)`.
- Debug function is similarly passed via context (key: `"debugFunc"`).
- Tests use interfaces (e.g., `catAliasesSvc`, `catIndicesSvc`) to mock ES service calls without a running cluster.

## Dependencies

- **github.com/olivere/elastic/v7** — Elasticsearch client
- **github.com/urfave/cli/v2** — CLI framework
- **github.com/BurntSushi/toml** — Config file parsing
- **github.com/cheggaaa/pb/v3** — Progress bar (bulk operations)
- **github.com/sergi/go-diff** — Diff output (search template render)
