# go-hr-core

Reusable HR domain library written in Go. Provides the canonical types and
use cases for any HR product (workers, vacancies, hiring) without binding to
a specific database, HTTP framework or messaging system.

`hases-api` and any future HR project (e.g. another company's RR.HH. portal)
can build on top of this module instead of duplicating the domain model.

## Install

```bash
go get github.com/nick130920/go-hr-core
```

## Packages

| Path | Responsibility |
|---|---|
| `worker` | `Worker` aggregate, status state-machine and `Service` use cases (hire, terminate, list). |
| `vacancy` | `Vacancy` and `Application` aggregates with publish/close transitions. |
| `hiring` | Cross-aggregate coordinator that hires a worker and closes the originating vacancy in one call. |

## Architecture

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

The library only ships the domain (entities + service + repository
interfaces). Adapters live in the consuming application:

- `hases-api/internal/adapters/persistence/pg_workers.go` implements
  `worker.Repository` against pgx.
- `hases-api/internal/adapters/http/handlers_worker.go` becomes a thin
  controller delegating to `worker.Service`.

## Quick start

```go
import (
    "context"
    "time"

    "github.com/nick130920/go-hr-core/hiring"
    "github.com/nick130920/go-hr-core/vacancy"
    "github.com/nick130920/go-hr-core/worker"
)

workerSvc := worker.NewService(workerRepo, time.Now)
vacancySvc := vacancy.NewService(vacancyRepo, time.Now)
coord := hiring.New(workerSvc, vacancySvc, eventBus.Publish, time.Now)

if err := coord.Hire(ctx, "hases", vacancyID, workerID, "Ana Pérez"); err != nil {
    return fmt.Errorf("hire: %w", err)
}
```

## License

Apache-2.0.
