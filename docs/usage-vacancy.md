# Vacancy service

The `vacancy` package owns vacancies and their applications.

## Aggregates

```go
type Vacancy struct {
    ID         uuid.UUID
    TenantID   string
    Title      string
    Department string
    State      State        // draft | published | closed
    OpenedAt   time.Time
    ClosedAt   *time.Time
}

type Application struct {
    ID          uuid.UUID
    VacancyID   uuid.UUID
    TenantID    string
    Candidate   string
    Email       string
    SubmittedAt time.Time
    Score       int
}
```

## State-machine

```
draft ──▶ published ──▶ closed
   \                       ▲
    \─────────────────────/
```

`Service.Publish` only succeeds from `draft`; `Service.Close` is allowed
from any state so that emergencies (legal hold, fraud) can take effect
immediately. The transition is rejected with `ErrInvalidTransition` when
the rule is violated.

## Service

```go
import (
    "time"

    "github.com/nick130920/go-hr-core/vacancy"
)

svc := vacancy.NewService(repo, time.Now)

if err := svc.Publish(ctx, "hases", vacancyID); err != nil {
    if errors.Is(err, vacancy.ErrInvalidTransition) {
        return fiber.ErrConflict
    }
    return err
}
```

## Repository contract

```go
type Repository interface {
    CreateVacancy(ctx context.Context, v *Vacancy) error
    PublishVacancy(ctx context.Context, tenantID string, id uuid.UUID, at time.Time) error
    CloseVacancy(ctx context.Context, tenantID string, id uuid.UUID, at time.Time) error
    GetVacancy(ctx context.Context, tenantID string, id uuid.UUID) (*Vacancy, error)
    CreateApplication(ctx context.Context, a *Application) error
    ListApplications(ctx context.Context, tenantID string, vacancyID uuid.UUID) ([]Application, error)
}
```
