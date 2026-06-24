package findappointment

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/ports"
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

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Appointment, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Appointment{}, recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil {
		return domain.Appointment{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.AppointmentByID(ctx, cmd.PsychologistID, cmd.ID)
}
