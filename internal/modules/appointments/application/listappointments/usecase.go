package listappointments

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	cfg  appshared.Config
	repo ports.Repository
}

func New(cfg appshared.Config, repo ports.Repository) *UseCase {
	return &UseCase{cfg: cfg, repo: repo}
}

type Command struct {
	PsychologistID uuid.UUID
	Pagination     appshared.CursorPagination
	PatientID      *uuid.UUID
	StartDate      *time.Time
	EndDate        *time.Time
	Statuses       []sharedv1.AppointmentStatus
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.AppointmentList, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.AppointmentList{}, recorderrors.ErrUnauthenticated
	}
	pagination, err := appshared.ResolveCursorPagination(cmd.Pagination, uc.cfg)
	if err != nil {
		return domain.AppointmentList{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.ListAppointments(ctx, domain.AppointmentListFilter{
		PsychologistID: cmd.PsychologistID,
		Pagination:     pagination,
		PatientID:      cmd.PatientID,
		StartDate:      cmd.StartDate,
		EndDate:        cmd.EndDate,
		Statuses:       cmd.Statuses,
	})
}
