package createattachment

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	repo    ports.Repository
	files   ports.FileManagementClient
	now     func() time.Time
	newUUID func() uuid.UUID
}

func New(repo ports.Repository, files ports.FileManagementClient, runtime appshared.Runtime) *UseCase {
	runtime = runtime.Normalize()
	return &UseCase{repo: repo, files: files, now: runtime.Now, newUUID: runtime.NewUUID}
}

type Command struct {
	PsychologistID uuid.UUID
	UserID         *uuid.UUID
	PatientID      uuid.UUID
	NoteID         *uuid.UUID
	OriginalName   string
	MimeType       string
	SizeBytes      int64
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.AttachmentCreateResult, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.AttachmentCreateResult{}, recorderrors.ErrUnauthenticated
	}
	if cmd.PatientID == uuid.Nil || strings.TrimSpace(cmd.OriginalName) == "" || strings.TrimSpace(cmd.MimeType) == "" || cmd.SizeBytes <= 0 {
		return domain.AttachmentCreateResult{}, recorderrors.ErrInvalidInput
	}
	exists, err := uc.repo.PatientExists(ctx, cmd.PsychologistID, cmd.PatientID)
	if err != nil {
		return domain.AttachmentCreateResult{}, err
	}
	if !exists {
		return domain.AttachmentCreateResult{}, recorderrors.ErrNotFound
	}
	if cmd.NoteID != nil {
		ok, err := uc.repo.NoteBelongsToPatient(ctx, cmd.PsychologistID, *cmd.NoteID, cmd.PatientID)
		if err != nil {
			return domain.AttachmentCreateResult{}, err
		}
		if !ok {
			return domain.AttachmentCreateResult{}, recorderrors.ErrInvalidInput
		}
	}

	fileID, uploadURL, expiresAt, err := uc.files.RequestUpload(ctx, strings.TrimSpace(cmd.OriginalName), strings.TrimSpace(cmd.MimeType), cmd.SizeBytes, false, "record")
	if err != nil {
		return domain.AttachmentCreateResult{}, err
	}

	attachmentID := uc.newUUID()
	attachment, err := uc.repo.CreateAttachment(ctx, domain.AttachmentCreate{
		ID:             attachmentID,
		PsychologistID: cmd.PsychologistID,
		UserID:         cmd.UserID,
		PatientID:      cmd.PatientID,
		NoteID:         cmd.NoteID,
		OriginalName:   strings.TrimSpace(cmd.OriginalName),
		MimeType:       strings.TrimSpace(cmd.MimeType),
		SizeBytes:      cmd.SizeBytes,
		Now:            uc.now().UTC(),
	}, fileID)
	if err != nil {
		return domain.AttachmentCreateResult{}, err
	}

	return domain.AttachmentCreateResult{
		ID:        attachment.ID,
		FileID:    attachment.FileID,
		UploadURL: uploadURL,
		ExpiresAt: expiresAt,
	}, nil
}
