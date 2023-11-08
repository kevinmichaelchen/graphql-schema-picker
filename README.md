# graphql-schema-picker

[![GoReportCard example](https://goreportcard.com/badge/github.com/kevinmichaelchen/graphql-schema-picker)](https://goreportcard.com/report/github.com/kevinmichaelchen/graphql-schema-picker)
[![version](https://img.shields.io/github/v/release/kevinmichaelchen/graphql-schema-picker?include_prereleases&label=latest&logo=ferrari)](https://github.com/kevinmichaelchen/graphql-schema-picker/releases/latest)
[![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability/kevinmichaelchen/graphql-schema-picker)](https://codeclimate.com/github/kevinmichaelchen/graphql-schema-picker)

<img width="600" src="https://github.com/kevinmichaelchen/graphql-schema-picker/assets/5129994/0b7c6707-a76f-4a49-9539-279969307fc8" />

A CLI for selectively pruning your GraphQL schemas.

The CLI accepts a Schema Definition Language file, and then selectively picks
(or filters out) certain elements.

## Motivation

This tool was born out of a desire to reuse Hasura's schema in upstream
microservices. In my case, it meant discarding the tens of thousands of lines
in the schema that was introspected from Hasura, and really only paying
attention to the few types I cared about.

## Example

For a realistic example of what a Hasura GraphQL schema looks like, check out
our example [**SDL file**][sdl-file] (Schema Definition Language).

[sdl-file]: ./examples/hasura.sdl.graphqls

## Getting Started

### Installing

Eventually, I may package this up in [pkgx][pkgx] and maybe even Homebrew (via
[Goreleaser][goreleaser-brew]).

[pkgx]: https://pkgx.sh/
[goreleaser-brew]: https://goreleaser.com/customization/homebrew/

#### With `go install`

```shell
go install github.com/kevinmichaelchen/graphql-schema-picker@latest
```

#### With Docker

```shell
docker pull ghcr.io/kevinmichaelchen/graphql-schema-picker
docker run --rm ghcr.io/kevinmichaelchen/graphql-schema-picker --help

docker run --rm \
  -v $(pwd)/examples:/examples \
  ghcr.io/kevinmichaelchen/graphql-schema-picker \
    --debug \
    pick \
      --output /examples/pruned.sdl.graphqls \
      --sdl-file /examples/hasura.sdl.graphqls \
      --definitions Aircrafts
```

## Similar Tools

- https://github.com/n1ru4l/graphql-public-schema-filter
- https://github.com/kesne/graphql-schema-subset
- https://github.com/xometry/graphql-code-generator-subset-plugin
- https://the-guild.dev/graphql/tools/docs/api/classes/wrap_src.pruneschema
- https://pothos-graphql.dev/docs/plugins/sub-graph

## Contributing

### Building

```shell
go run cmd/graphql-schema-picker/main.go \
  --debug \
  pick \
    --output examples/pruned.sdl.graphqls \
    --sdl-file examples/hasura.sdl.graphqls \
    --definitions Aircrafts
```

### Releasing

Follow [Conventional Commits][conventional-commits] and SemVer releases should
happen automatically via GitHub Actions. 

[conventional-commits]: https://www.conventionalcommits.org/en/v1.0.0/

## Tasks

### build

Builds the Go program into a local binary.

```shell
pkgx goreleaser build --clean --single-target
```
