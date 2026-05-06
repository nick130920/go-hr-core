# Hiring coordinator

The `hiring` package wires `worker` and `vacancy` together so that callers
do not orchestrate two services from inside a handler.

## Why a coordinator?

Hiring is the only HR operation that crosses two aggregates: a worker
becomes active **and** a vacancy gets closed. Without the coordinator the
caller would:

- Open a transaction.
- Update the worker.
- Update the vacancy.
- Publish a domain event.
- Close the transaction.

The coordinator captures that flow once and exposes a single method to
consumers.

## Usage

```go
import (
    "context"
    "time"

    "github.com/google/uuid"

    "github.com/nick130920/go-hr-core/hiring"
    "github.com/nick130920/go-hr-core/vacancy"
    "github.com/nick130920/go-hr-core/worker"
)

workerSvc  := worker.NewService(workerRepo, time.Now)
vacancySvc := vacancy.NewService(vacancyRepo, time.Now)
publish    := func(ctx context.Context, e hiring.Event) error {
    return bus.Send(ctx, "hr.hire.completed", e)
}

coord := hiring.New(workerSvc, vacancySvc, publish, time.Now)

err := coord.Hire(ctx, "hases", vacancyID, workerID, "Ana Pérez")
if err != nil {
    return fmt.Errorf("hire: %w", err)
}
```

## Event payload

The published `Event` is the audit record of the hire:

```go
type Event struct {
    TenantID   string
    WorkerID   uuid.UUID
    VacancyID  uuid.UUID
    Candidate  string
    HappenedAt time.Time
}
```

Downstream systems (induction, payroll, IT provisioning) consume this
event to start their own pipelines without having to poll the database.
