package listpatients

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/ports"
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
	PsychologistID          uuid.UUID
	Pagination              appshared.CursorPagination
	WithPendingAppointments bool
	WithNoAppointments      bool
	EverHadAppointments     bool
	SearchQuery             string
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.PatientList, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.PatientList{}, recorderrors.ErrUnauthenticated
	}
	pagination, err := appshared.ResolveCursorPagination(cmd.Pagination, uc.cfg)
	if err != nil {
		return domain.PatientList{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.ListPatients(ctx, domain.PatientListFilter{
		PsychologistID:          cmd.PsychologistID,
		Pagination:              pagination,
		SearchQuery:             strings.TrimSpace(cmd.SearchQuery),
		WithPendingAppointments: cmd.WithPendingAppointments,
		WithNoAppointments:      cmd.WithNoAppointments,
		EverHadAppointments:     cmd.EverHadAppointments,
	})
}
