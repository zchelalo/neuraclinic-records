package updatepatientidentification

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
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
	FirstName      *string
	MiddleName     *string
	FirstLastName  *string
	SecondLastName *string
	BirthDate      *time.Time
	Sex            *recordv1.Sex
	BirthCountry   *string
	BirthProvince  *string
	BirthCity      *string
	Occupation     *string
	MaritalStatus  *recordv1.MaritalStatus
	Religion       *string
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Patient, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Patient{}, recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil {
		return domain.Patient{}, recorderrors.ErrInvalidInput
	}
	if cmd.Sex != nil && *cmd.Sex == recordv1.Sex_SEX_UNSPECIFIED {
		return domain.Patient{}, recorderrors.ErrInvalidInput
	}
	if cmd.MaritalStatus != nil && *cmd.MaritalStatus == recordv1.MaritalStatus_MARITAL_STATUS_UNSPECIFIED {
		return domain.Patient{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.UpdatePatientIdentification(ctx, domain.PatientIdentificationUpdate{
		ID:             cmd.ID,
		PsychologistID: cmd.PsychologistID,
		FirstName:      trimPtr(cmd.FirstName),
		MiddleName:     trimPtr(cmd.MiddleName),
		FirstLastName:  trimPtr(cmd.FirstLastName),
		SecondLastName: trimPtr(cmd.SecondLastName),
		BirthDate:      cmd.BirthDate,
		Sex:            cmd.Sex,
		BirthCountry:   trimPtr(cmd.BirthCountry),
		BirthProvince:  trimPtr(cmd.BirthProvince),
		BirthCity:      trimPtr(cmd.BirthCity),
		Occupation:     trimPtr(cmd.Occupation),
		MaritalStatus:  cmd.MaritalStatus,
		Religion:       trimPtr(cmd.Religion),
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
