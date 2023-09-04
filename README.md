# graphql-schema-picker

A CLI for selectively pruning your GraphQL schemas.

The CLI accepts a Schema Definition Language file, and then selectively picks 
(or filters out) certain elements.

## Motivation

This tool was born out of a desire to reuse Hasura's schema in upstream 
microservices. In my case, it meant discarding the tens of thousands of lines
in the schema that was introspected from Hasura, and really only paying 
attention to the few types I cared about.

## Example

```graphql
schema {
  query: query_root
  mutation: mutation_root
  subscription: subscription_root
}

type query_root {
}
```

## Getting Started

### Installing

### Usage

```shell
go run cmd/graphql-schema-picker/main.go \
  --debug \
  pick \
    --sdl-file examples/hasura.sdl.graphqls \
    --definitions Aircrafts
```