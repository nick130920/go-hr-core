// Package worker models people that have been hired (or are being onboarded)
// in an organization. The package is deliberately framework-agnostic: it
// only depends on the standard library and `github.com/google/uuid` so it
// can be reused by any HR-flavored service.
package worker

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Status enumerates the worker lifecycle stages used across the platform.
type Status string

const (
	StatusInvited    Status = "invited"
	StatusInducting  Status = "inducting"
	StatusActive     Status = "active"
	StatusOnLeave    Status = "on_leave"
	StatusTerminated Status = "terminated"
)

// Worker is the canonical representation of an employee in the HR domain.
// Companies extend it with custom fields via embedding; the core fields
// below are the minimum every HR product needs.
type Worker struct {
	ID        uuid.UUID
	TenantID  string
	FirstName string
	LastName  string
	Email     string
	NationalID string
	Status    Status
	HiredAt   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repository is the persistence contract. Concrete adapters (Postgres, Mongo,
// in-memory for tests) live in the consuming service.
type Repository interface {
	Create(ctx context.Context, w *Worker) error
	Update(ctx context.Context, w *Worker) error
	GetByID(ctx context.Context, tenantID string, id uuid.UUID) (*Worker, error)
	GetByEmail(ctx context.Context, tenantID, email string) (*Worker, error)
	List(ctx context.Context, tenantID string, q ListQuery) ([]Worker, error)
}

// ListQuery captures the supported filters and pagination knobs.
type ListQuery struct {
	Status []Status
	Search string
	Limit  int
	Offset int
}

// ErrNotFound is returned by repositories when a worker does not exist.
var ErrNotFound = errors.New("worker: not found")

// Service exposes the use cases consumers should call. It hides repository
// details and centralizes invariants such as "an active worker must have a
// hire date".
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService wires the dependencies. now defaults to time.Now when nil.
func NewService(repo Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{repo: repo, now: now}
}

// Hire transitions a worker into the active state and persists the hire
// date. It returns ErrNotFound when the worker cannot be located.
func (s *Service) Hire(ctx context.Context, tenantID string, id uuid.UUID) (*Worker, error) {
	w, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	w.Status = StatusActive
	w.HiredAt = &now
	w.UpdatedAt = now
	if err := s.repo.Update(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

// Terminate moves a worker to the terminated state.
func (s *Service) Terminate(ctx context.Context, tenantID string, id uuid.UUID) (*Worker, error) {
	w, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	w.Status = StatusTerminated
	w.UpdatedAt = s.now().UTC()
	if err := s.repo.Update(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}
