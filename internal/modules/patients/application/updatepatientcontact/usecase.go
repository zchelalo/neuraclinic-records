package updatepatientcontact

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/ports"
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
	Phone          *string
	Email          *string
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Patient, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Patient{}, recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil {
		return domain.Patient{}, recorderrors.ErrInvalidInput
	}
	email := trimPtr(cmd.Email)
	if email != nil {
		lower := strings.ToLower(*email)
		email = &lower
	}
	return uc.repo.UpdatePatientContact(ctx, domain.PatientContactUpdate{
		ID:             cmd.ID,
		PsychologistID: cmd.PsychologistID,
		Phone:          trimPtr(cmd.Phone),
		Email:          email,
		Now:            uc.now().UTC(),
	})
}

func trimPtr(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}
