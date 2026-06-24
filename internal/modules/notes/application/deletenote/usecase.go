package deletenote

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	repo ports.Repository
	now  func() time.Time
}

func New(repo ports.Repository, runtime appshared.Runtime) *UseCase {
	runtime = runtime.Normalize()
	return &UseCase{repo: repo, now: runtime.Now}
}

type Command struct {
	PsychologistID uuid.UUID
	ID             uuid.UUID
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) error {
	if cmd.PsychologistID == uuid.Nil {
		return recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil {
		return recorderrors.ErrInvalidInput
	}
	return uc.repo.DeleteNote(ctx, cmd.PsychologistID, cmd.ID, uc.now().UTC())
}
