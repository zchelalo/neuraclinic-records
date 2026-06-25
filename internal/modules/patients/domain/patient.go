package domain

import (
	"time"

	"github.com/google/uuid"
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
)

type Address struct {
	ID           uuid.UUID
	Country      string
	Province     string
	City         string
	PostalCode   string
	Neighborhood string
	Street       string
	StreetNumber string
	UnitNumber   *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type Patient struct {
	ID             uuid.UUID
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
	Address        Address
	PsychologistID uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type PatientSummary struct {
	ID             uuid.UUID
	FirstName      string
	MiddleName     *string
	FirstLastName  string
	SecondLastName *string
	BirthDate      time.Time
	Email          string
	Phone          string
}

type PatientCreate struct {
	ID             uuid.UUID
	AddressID      uuid.UUID
	FamiliogramID  uuid.UUID
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
	Now            time.Time
}

type PatientListFilter struct {
	PsychologistID          uuid.UUID
	Pagination              appshared.ResolvedCursorPagination
	SearchQuery             string
	WithPendingAppointments bool
	WithNoAppointments      bool
	EverHadAppointments     bool
}

type PatientList struct {
	Patients []PatientSummary
	Meta     appshared.CursorMeta
}

type PatientIdentificationUpdate struct {
	ID             uuid.UUID
	PsychologistID uuid.UUID
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
	Now            time.Time
}

type PatientContactUpdate struct {
	ID             uuid.UUID
	PsychologistID uuid.UUID
	Phone          *string
	Email          *string
	Now            time.Time
}

type AddressUpdate struct {
	PatientID      uuid.UUID
	PsychologistID uuid.UUID
	Country        *string
	Province       *string
	City           *string
	PostalCode     *string
	Neighborhood   *string
	Street         *string
	StreetNumber   *string
	UnitNumber     *string
	Now            time.Time
}
