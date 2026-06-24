package domain

import (
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
)

type Appointment struct {
	ID                           uuid.UUID
	StartTime                    time.Time
	EndTime                      time.Time
	Reason                       string
	Status                       sharedv1.AppointmentStatus
	PatientID                    uuid.UUID
	CancelledByUserID            *uuid.UUID
	RescheduledFromAppointmentID *uuid.UUID
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}

type AppointmentCreate struct {
	ID                           uuid.UUID
	PsychologistID               uuid.UUID
	PatientID                    uuid.UUID
	StartTime                    time.Time
	EndTime                      time.Time
	Reason                       string
	Status                       sharedv1.AppointmentStatus
	CancelledByUserID            *uuid.UUID
	RescheduledFromAppointmentID *uuid.UUID
	Now                          time.Time
}

type AppointmentListFilter struct {
	PsychologistID uuid.UUID
	Pagination     appshared.ResolvedCursorPagination
	PatientID      *uuid.UUID
	StartDate      *time.Time
	EndDate        *time.Time
	Statuses       []sharedv1.AppointmentStatus
}

type AppointmentList struct {
	Appointments []Appointment
	Meta         appshared.CursorMeta
}
