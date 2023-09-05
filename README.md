# graphql-schema-picker

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

Eventually, I may package this up in Tea and maybe even Homebrew (via 
[Goreleaser][goreleaser-brew]).

For the time being, it should be installable with Go:

```shell
go install github.com/kevinmichaelchen/graphql-schema-picker@latest
```

[goreleaser-brew]: https://goreleaser.com/customization/homebrew/

### Usage

```shell
go run cmd/graphql-schema-picker/main.go \
  --debug \
  pick \
    --sdl-file examples/hasura.sdl.graphqls \
    --definitions Aircrafts
```

## Similar Tools

- https://github.com/n1ru4l/graphql-public-schema-filter
- https://github.com/kesne/graphql-schema-subset
- https://github.com/xometry/graphql-code-generator-subset-plugin
- https://the-guild.dev/graphql/tools/docs/api/classes/wrap_src.pruneschema
- https://pothos-graphql.dev/docs/plugins/sub-graph