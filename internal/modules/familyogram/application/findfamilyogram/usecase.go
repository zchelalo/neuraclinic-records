package findfamilyogram

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/ports"
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
	PatientID      uuid.UUID
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Familyogram, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Familyogram{}, recorderrors.ErrUnauthenticated
	}
	if cmd.PatientID == uuid.Nil {
		return domain.Familyogram{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.FamilyogramByPatientID(ctx, cmd.PsychologistID, cmd.PatientID)
}
