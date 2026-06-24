package createpatient

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
	repo    ports.Repository
	now     func() time.Time
	newUUID func() uuid.UUID
}

func New(repo ports.Repository, runtime appshared.Runtime) *UseCase {
	runtime = runtime.Normalize()
	return &UseCase{repo: repo, now: runtime.Now, newUUID: runtime.NewUUID}
}

type Command struct {
	PsychologistID uuid.UUID
	FirstName      string
	MiddleName     *string
	FirstLastName  string
	SecondLastName *string
	BirthDate      time.Time
	BirthCountry   string
	BirthProvince  string
	BirthCity      string
	Sex            recordv1.Sex
	MaritalStatus  recordv1.MaritalStatus
	Occupation     *string
	Religion       *string
	Phone          string
	Email          string
	Country        string
	Province       string
	City           string
	PostalCode     string
	Neighborhood   string
	Street         string
	StreetNumber   string
	UnitNumber     *string
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.PatientSummary, error) {
	if cmd.PsychologistID == uuid.Nil ||
		blank(cmd.FirstName, cmd.FirstLastName, cmd.BirthCountry, cmd.BirthProvince, cmd.BirthCity, cmd.Phone, cmd.Email, cmd.Country, cmd.Province, cmd.City, cmd.PostalCode, cmd.Neighborhood, cmd.Street, cmd.StreetNumber) ||
		cmd.BirthDate.IsZero() ||
		cmd.Sex == recordv1.Sex_SEX_UNSPECIFIED ||
		cmd.MaritalStatus == recordv1.MaritalStatus_MARITAL_STATUS_UNSPECIFIED {
		return domain.PatientSummary{}, recorderrors.ErrInvalidInput
	}

	now := uc.now().UTC()
	return uc.repo.CreatePatient(ctx, domain.PatientCreate{
		ID:             uc.newUUID(),
		AddressID:      uc.newUUID(),
		FamilyogramID:  uc.newUUID(),
		PsychologistID: cmd.PsychologistID,
		FirstName:      strings.TrimSpace(cmd.FirstName),
		MiddleName:     trimPtr(cmd.MiddleName),
		FirstLastName:  strings.TrimSpace(cmd.FirstLastName),
		SecondLastName: trimPtr(cmd.SecondLastName),
		BirthDate:      cmd.BirthDate,
		BirthCountry:   strings.TrimSpace(cmd.BirthCountry),
		BirthProvince:  strings.TrimSpace(cmd.BirthProvince),
		BirthCity:      strings.TrimSpace(cmd.BirthCity),
		Sex:            cmd.Sex,
		MaritalStatus:  cmd.MaritalStatus,
		Occupation:     trimPtr(cmd.Occupation),
		Religion:       trimPtr(cmd.Religion),
		Phone:          strings.TrimSpace(cmd.Phone),
		Email:          strings.TrimSpace(strings.ToLower(cmd.Email)),
		Country:        strings.TrimSpace(cmd.Country),
		Province:       strings.TrimSpace(cmd.Province),
		City:           strings.TrimSpace(cmd.City),
		PostalCode:     strings.TrimSpace(cmd.PostalCode),
		Neighborhood:   strings.TrimSpace(cmd.Neighborhood),
		Street:         strings.TrimSpace(cmd.Street),
		StreetNumber:   strings.TrimSpace(cmd.StreetNumber),
		UnitNumber:     trimPtr(cmd.UnitNumber),
		Now:            now,
	})
}

func blank(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}

func trimPtr(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}
