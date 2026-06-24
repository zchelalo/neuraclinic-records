package updatefamilyogram

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"google.golang.org/protobuf/types/known/structpb"
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
	Data           *structpb.Struct
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Familyogram, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Familyogram{}, recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil || cmd.Data == nil {
		return domain.Familyogram{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.UpdateFamilyogram(ctx, cmd.PsychologistID, cmd.ID, cmd.Data, uc.now().UTC())
}
