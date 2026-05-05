// Package hiring orchestrates worker creation and vacancy follow-up so
// callers do not have to wire two repositories themselves.
package hiring

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nick130920/go-hr-core/vacancy"
	"github.com/nick130920/go-hr-core/worker"
)

// Coordinator runs cross-package operations: hiring a candidate, closing the
// vacancy that produced the hire, and emitting an event for downstream
// systems (induction, payroll, ...).
type Coordinator struct {
	workers   *worker.Service
	vacancies *vacancy.Service
	publish   func(ctx context.Context, e Event) error
	now       func() time.Time
}

// Event is the audit record emitted when a candidate is hired.
type Event struct {
	TenantID    string
	WorkerID    uuid.UUID
	VacancyID   uuid.UUID
	Candidate   string
	HappenedAt  time.Time
}

// New wires the coordinator. publish may be nil if no event bus is needed.
func New(w *worker.Service, v *vacancy.Service, publish func(ctx context.Context, e Event) error, now func() time.Time) *Coordinator {
	if now == nil {
		now = time.Now
	}
	return &Coordinator{workers: w, vacancies: v, publish: publish, now: now}
}

// Hire transitions the worker to active, closes the vacancy and emits an
// event when configured.
func (c *Coordinator) Hire(ctx context.Context, tenantID string, vacancyID, workerID uuid.UUID, candidate string) error {
	if _, err := c.workers.Hire(ctx, tenantID, workerID); err != nil {
		return err
	}
	if err := c.vacancies.Close(ctx, tenantID, vacancyID); err != nil {
		return err
	}
	if c.publish != nil {
		return c.publish(ctx, Event{
			TenantID:   tenantID,
			WorkerID:   workerID,
			VacancyID:  vacancyID,
			Candidate:  candidate,
			HappenedAt: c.now().UTC(),
		})
	}
	return nil
}
