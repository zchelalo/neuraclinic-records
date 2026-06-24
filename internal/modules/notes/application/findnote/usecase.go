package findnote

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/ports"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	repo ports.Repository
}

func New(repo ports.Repository) *UseCase {
	return &UseCase{repo: repo}
}

type Command struct {
	PsychologistID uuid.UUID
	ID             uuid.UUID
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Note, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Note{}, recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil {
		return domain.Note{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.NoteByID(ctx, cmd.PsychologistID, cmd.ID)
}
