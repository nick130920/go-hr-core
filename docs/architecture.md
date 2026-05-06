# Architecture

`go-hr-core` follows a Domain-Driven Design layout: every package owns one
aggregate and exposes a `Service` plus a `Repository` interface. The
package boundaries match the natural transactional boundaries.

## Module layout

```
+-------------------+      +--------------------+
|  HTTP / Worker    | ---> |  hiring.Coordinator|
|  (chi handlers)   |      +--------------------+
+-------------------+         |              |
                              v              v
                       worker.Service   vacancy.Service
                              |              |
                              v              v
                     worker.Repository  vacancy.Repository
                              \            /
                               \          /
                                Postgres / Mongo / in-mem
```

## Layered responsibilities

| Layer | What lives here | Owned by |
|---|---|---|
| **Domain** (this lib) | Aggregates, value objects, services, repository interfaces. No I/O, no framework. | `go-hr-core` |
| **Adapters** | Concrete repositories (pgx, GORM), HTTP handlers, queue consumers. | The consuming app, e.g. `hases-api` |
| **Composition root** | Wiring: builds the repositories, injects them into services, mounts routes. | `cmd/api/main.go` of the consuming app |

## Why `Service` is in this library

In a strict hexagonal layout, services would live in the application layer
of each consumer. We deliberately put them here because the **rules** of
"a worker can only be hired once", "a vacancy moves draft → published →
closed", "the hire date is recorded at the moment of hire" are universal
across HR products. Pushing the logic into each consumer would re-introduce
the duplication this library exists to eliminate.

## Extending the model

Companies extend the canonical types via embedding:

```go
type HasesWorker struct {
    worker.Worker
    Department string
    NationalID string
}
```

Custom fields stay in the consumer's database schema; the library does not
need to know about them. Repositories return the canonical type, and the
consumer up-casts when needed.
