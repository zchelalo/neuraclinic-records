package ports

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
)

type Repository interface {
	CreateAttachment(ctx context.Context, attachment domain.AttachmentCreate, fileID uuid.UUID) (domain.Attachment, error)
	AttachmentByID(ctx context.Context, psychologistID, id uuid.UUID) (domain.Attachment, error)
	ListAttachments(ctx context.Context, filter domain.AttachmentListFilter) (domain.AttachmentList, error)
	UpdateAttachmentUploadStatusByFileID(ctx context.Context, fileID uuid.UUID, status sharedv1.FileStatus, now time.Time) (domain.Attachment, error)
	DeleteAttachment(ctx context.Context, psychologistID, id uuid.UUID, now time.Time) error
	PatientExists(ctx context.Context, psychologistID, id uuid.UUID) (bool, error)
	NoteBelongsToPatient(ctx context.Context, psychologistID, noteID, patientID uuid.UUID) (bool, error)
}

var ErrInvalidEvent = errors.New("invalid event")
