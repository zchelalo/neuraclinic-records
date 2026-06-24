package listnotes

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/ports"
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
	PsychologistID            uuid.UUID
	PatientID                 uuid.UUID
	Pagination                appshared.CursorPagination
	StartDate                 *time.Time
	EndDate                   *time.Time
	WithAppointmentAssociated bool
	WithFilesAssociated       bool
	SearchQuery               string
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.NoteList, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.NoteList{}, recorderrors.ErrUnauthenticated
	}
	if cmd.PatientID == uuid.Nil {
		return domain.NoteList{}, recorderrors.ErrInvalidInput
	}
	exists, err := uc.repo.PatientExists(ctx, cmd.PsychologistID, cmd.PatientID)
	if err != nil {
		return domain.NoteList{}, err
	}
	if !exists {
		return domain.NoteList{}, recorderrors.ErrNotFound
	}
	pagination, err := appshared.ResolveCursorPagination(cmd.Pagination, uc.cfg)
	if err != nil {
		return domain.NoteList{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.ListNotes(ctx, domain.NoteListFilter{
		PsychologistID:            cmd.PsychologistID,
		PatientID:                 cmd.PatientID,
		Pagination:                pagination,
		StartDate:                 cmd.StartDate,
		EndDate:                   cmd.EndDate,
		WithAppointmentAssociated: cmd.WithAppointmentAssociated,
		WithFilesAssociated:       cmd.WithFilesAssociated,
		SearchQuery:               strings.TrimSpace(cmd.SearchQuery),
	})
}
