A small example to generate realistic Hasura schemas.

## Entity Relationship Diagram

<img src="./diagrams/erd.svg" width="500" />

## Tasks

### start

Start Hasura

```shell
docker compose up --detach
```

### console

Launch the Hasura Console UI

```shell
hasura --project hasura \
  console
```

### migrate_format

Format migrations.

```shell
sleek -n \
  --indent-spaces 2 \
  --uppercase \
  hasura/migrations/**/*.sql
```

### migrate_status

List migrations

```shell
hasura --project hasura \
  migrate \
    --database-name default \
    status
```

### migrate_squash

Squash all migrations

```shell
hasura --project hasura \
  migrate \
    --database-name default \
    squash --delete-source \
    --from 1693763548205 \
    --name init
```

### introspect

Introspect the Hasura schema — i.e., the Schema Definition Language (SDL) file.

```shell
npx graphqurl -- \
  http://localhost:8080/v1/graphql \
  --introspect > hasura.sdl.graphqls
```
