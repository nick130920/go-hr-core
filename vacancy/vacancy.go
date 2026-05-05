// Package vacancy models open positions and candidate applications.
package vacancy

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// State is the lifecycle stage of a vacancy.
type State string

const (
	StateDraft     State = "draft"
	StatePublished State = "published"
	StateClosed    State = "closed"
)

// Vacancy is an open position the company wants to fill.
type Vacancy struct {
	ID         uuid.UUID
	TenantID   string
	Title      string
	Department string
	State      State
	OpenedAt   time.Time
	ClosedAt   *time.Time
}

// Application is a candidate response to a vacancy.
type Application struct {
	ID         uuid.UUID
	VacancyID  uuid.UUID
	TenantID   string
	Candidate  string
	Email      string
	SubmittedAt time.Time
	Score      int
}

// Repository persists vacancies and their applications.
type Repository interface {
	CreateVacancy(ctx context.Context, v *Vacancy) error
	PublishVacancy(ctx context.Context, tenantID string, id uuid.UUID, at time.Time) error
	CloseVacancy(ctx context.Context, tenantID string, id uuid.UUID, at time.Time) error
	GetVacancy(ctx context.Context, tenantID string, id uuid.UUID) (*Vacancy, error)

	CreateApplication(ctx context.Context, a *Application) error
	ListApplications(ctx context.Context, tenantID string, vacancyID uuid.UUID) ([]Application, error)
}

// ErrInvalidTransition is returned by Service when a vacancy can't move to
// the requested state.
var ErrInvalidTransition = errors.New("vacancy: invalid transition")

// Service centralizes the state-machine rules for vacancies.
type Service struct {
	repo Repository
	now  func() time.Time
}

// NewService wires the dependencies; now defaults to time.Now.
func NewService(r Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{repo: r, now: now}
}

// Publish makes a draft vacancy public.
func (s *Service) Publish(ctx context.Context, tenantID string, id uuid.UUID) error {
	v, err := s.repo.GetVacancy(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if v.State != StateDraft {
		return ErrInvalidTransition
	}
	return s.repo.PublishVacancy(ctx, tenantID, id, s.now().UTC())
}

// Close ends a vacancy regardless of its current state.
func (s *Service) Close(ctx context.Context, tenantID string, id uuid.UUID) error {
	return s.repo.CloseVacancy(ctx, tenantID, id, s.now().UTC())
}
