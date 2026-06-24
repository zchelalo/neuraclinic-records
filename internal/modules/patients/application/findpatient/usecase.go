package findpatient

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/ports"
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

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Patient, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Patient{}, recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil {
		return domain.Patient{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.PatientByID(ctx, cmd.PsychologistID, cmd.ID)
}
