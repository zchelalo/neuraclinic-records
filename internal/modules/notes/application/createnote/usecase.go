package createnote

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/ports"
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
	PatientID      uuid.UUID
	AppointmentID  *uuid.UUID
	Title          *string
	ContentHTML    string
	ContentText    string
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Note, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Note{}, recorderrors.ErrUnauthenticated
	}
	if cmd.PatientID == uuid.Nil || strings.TrimSpace(cmd.ContentHTML) == "" || strings.TrimSpace(cmd.ContentText) == "" {
		return domain.Note{}, recorderrors.ErrInvalidInput
	}
	exists, err := uc.repo.PatientExists(ctx, cmd.PsychologistID, cmd.PatientID)
	if err != nil {
		return domain.Note{}, err
	}
	if !exists {
		return domain.Note{}, recorderrors.ErrNotFound
	}
	if cmd.AppointmentID != nil {
		ok, err := uc.repo.AppointmentBelongsToPatient(ctx, cmd.PsychologistID, *cmd.AppointmentID, cmd.PatientID)
		if err != nil {
			return domain.Note{}, err
		}
		if !ok {
			return domain.Note{}, recorderrors.ErrInvalidInput
		}
	}
	return uc.repo.CreateNote(ctx, domain.NoteCreate{
		ID:             uc.newUUID(),
		PsychologistID: cmd.PsychologistID,
		PatientID:      cmd.PatientID,
		AppointmentID:  cmd.AppointmentID,
		Title:          trimPtr(cmd.Title),
		ContentHTML:    cmd.ContentHTML,
		ContentText:    cmd.ContentText,
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
