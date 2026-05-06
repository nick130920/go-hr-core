# Worker service

The `worker` package owns the employee aggregate and its lifecycle.

## Aggregate

```go
type Worker struct {
    ID         uuid.UUID
    TenantID   string
    FirstName  string
    LastName   string
    Email      string
    NationalID string
    Status     Status
    HiredAt    *time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

## Status state-machine

```
invited ──▶ inducting ──▶ active ──▶ on_leave
                              \         /
                               ▼       ▼
                            terminated
```

Transitions are enforced by the service, not by the domain type, so the
canonical `Worker` stays free of imperative branches.

## Service

```go
import (
    "time"

    "github.com/nick130920/go-hr-core/worker"
)

repo := newPostgresWorkerRepo(pool) // implements worker.Repository
svc  := worker.NewService(repo, time.Now)

w, err := svc.Hire(ctx, "hases", workerID)
if err != nil {
    return fmt.Errorf("hire: %w", err)
}
```

`Service.Hire` records the hire date and updates the status atomically;
`Service.Terminate` is the symmetric operation.

## Repository contract

```go
type Repository interface {
    Create(ctx context.Context, w *Worker) error
    Update(ctx context.Context, w *Worker) error
    GetByID(ctx context.Context, tenantID string, id uuid.UUID) (*Worker, error)
    GetByEmail(ctx context.Context, tenantID, email string) (*Worker, error)
    List(ctx context.Context, tenantID string, q ListQuery) ([]Worker, error)
}
```

Implement it once per persistence backend in your consuming app.
