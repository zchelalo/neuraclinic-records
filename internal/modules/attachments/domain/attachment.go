package domain

import (
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
)

type Attachment struct {
	ID           uuid.UUID
	FileID       uuid.UUID
	OriginalName string
	MimeType     string
	UploadStatus sharedv1.FileStatus
	DownloadURL  *string
	ExpiresAt    *time.Time
	PatientID    uuid.UUID
	NoteID       *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type AttachmentCreate struct {
	ID             uuid.UUID
	PsychologistID uuid.UUID
	UserID         *uuid.UUID
	PatientID      uuid.UUID
	NoteID         *uuid.UUID
	OriginalName   string
	MimeType       string
	SizeBytes      int64
	Now            time.Time
}

type AttachmentCreateResult struct {
	ID        uuid.UUID
	FileID    uuid.UUID
	UploadURL string
	ExpiresAt time.Time
}

type AttachmentListFilter struct {
	PsychologistID uuid.UUID
	PatientID      uuid.UUID
	NoteID         *uuid.UUID
	Pagination     appshared.ResolvedCursorPagination
}

type AttachmentList struct {
	Attachments []Attachment
	Meta        appshared.CursorMeta
}
