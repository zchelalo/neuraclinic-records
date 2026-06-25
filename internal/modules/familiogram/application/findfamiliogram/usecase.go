package findfamiliogram

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/ports"
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

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Familiogram, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Familiogram{}, recorderrors.ErrUnauthenticated
	}
	if cmd.PatientID == uuid.Nil {
		return domain.Familiogram{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.FamiliogramByPatientID(ctx, cmd.PsychologistID, cmd.PatientID)
}
