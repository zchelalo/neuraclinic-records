package updatepatientaddress

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
	Country        *string
	Province       *string
	City           *string
	PostalCode     *string
	Neighborhood   *string
	Street         *string
	StreetNumber   *string
	UnitNumber     *string
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Patient, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Patient{}, recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil {
		return domain.Patient{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.UpdatePatientAddress(ctx, domain.AddressUpdate{
		PatientID:      cmd.ID,
		PsychologistID: cmd.PsychologistID,
		Country:        trimPtr(cmd.Country),
		Province:       trimPtr(cmd.Province),
		City:           trimPtr(cmd.City),
		PostalCode:     trimPtr(cmd.PostalCode),
		Neighborhood:   trimPtr(cmd.Neighborhood),
		Street:         trimPtr(cmd.Street),
		StreetNumber:   trimPtr(cmd.StreetNumber),
		UnitNumber:     trimPtr(cmd.UnitNumber),
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
