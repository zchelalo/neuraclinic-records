package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/domain"
)

type Repository interface {
	CreateNote(ctx context.Context, note domain.NoteCreate) (domain.Note, error)
	NoteByID(ctx context.Context, psychologistID, id uuid.UUID) (domain.Note, error)
	ListNotes(ctx context.Context, filter domain.NoteListFilter) (domain.NoteList, error)
	UpdateNote(ctx context.Context, update domain.NoteUpdate) (domain.Note, error)
	DeleteNote(ctx context.Context, psychologistID, id uuid.UUID, now time.Time) error
	NoteBelongsToPatient(ctx context.Context, psychologistID, noteID, patientID uuid.UUID) (bool, error)
	PatientExists(ctx context.Context, psychologistID, id uuid.UUID) (bool, error)
	AppointmentBelongsToPatient(ctx context.Context, psychologistID, appointmentID, patientID uuid.UUID) (bool, error)
}
