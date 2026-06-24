package domain

import (
	"time"

	"github.com/google/uuid"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
)

type Note struct {
	ID            uuid.UUID
	PatientID     uuid.UUID
	AppointmentID *uuid.UUID
	Title         *string
	ContentHTML   string
	ContentText   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type NoteSummary struct {
	ID            uuid.UUID
	PatientID     uuid.UUID
	AppointmentID *uuid.UUID
	Title         *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type NoteCreate struct {
	ID             uuid.UUID
	PsychologistID uuid.UUID
	PatientID      uuid.UUID
	AppointmentID  *uuid.UUID
	Title          *string
	ContentHTML    string
	ContentText    string
	Now            time.Time
}

type NoteListFilter struct {
	PsychologistID            uuid.UUID
	PatientID                 uuid.UUID
	Pagination                appshared.ResolvedCursorPagination
	StartDate                 *time.Time
	EndDate                   *time.Time
	WithAppointmentAssociated bool
	WithFilesAssociated       bool
	SearchQuery               string
}

type NoteList struct {
	Notes []NoteSummary
	Meta  appshared.CursorMeta
}

type NoteUpdate struct {
	ID             uuid.UUID
	PsychologistID uuid.UUID
	AppointmentID  *uuid.UUID
	Title          *string
	ContentHTML    *string
	ContentText    *string
	Now            time.Time
}
