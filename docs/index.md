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
| [`worker`](usage-worker.md) | `Worker` aggregate, status state-machine and `Service` use cases (hire, terminate, list). |
| [`vacancy`](usage-vacancy.md) | `Vacancy` and `Application` aggregates with publish/close transitions. |
| [`hiring`](usage-hiring.md) | Cross-aggregate coordinator that hires a worker and closes the originating vacancy in one call. |

## Why this library exists

Every company that runs an HR product re-implements the same primitives:
employees with a status state-machine, vacancies with publish/close, an
application pipeline, a hire flow. By shipping these as a library:

- The domain model is **the same across companies**, so reports, audits and
  integrations are easier to standardize.
- Persistence and transport remain **per-product concerns**, keeping the
  library framework-agnostic.
- New companies can onboard with **a thin adapter layer** instead of
  re-writing the domain.
