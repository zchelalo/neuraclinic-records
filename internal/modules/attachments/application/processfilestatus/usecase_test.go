package processfilestatus

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	recordapp "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

func TestExecuteUpdatesAttachmentByFileID(t *testing.T) {
	ctx := context.Background()
	fileID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	occurredAt := time.Date(2026, 6, 25, 15, 0, 0, 0, time.UTC)
	repo := &processRepo{
		attachment: domain.Attachment{
			ID:           uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			FileID:       fileID,
			UploadStatus: sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
		},
	}

	uc := New(repo, recordapp.Runtime{})
	result, err := uc.Execute(ctx, Command{
		FileID:     fileID,
		Status:     sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
		OccurredAt: &occurredAt,
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if repo.lastFileID != fileID {
		t.Fatalf("expected file id %s, got %s", fileID, repo.lastFileID)
	}
	if repo.lastStatus != sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
		t.Fatalf("expected status available, got %s", repo.lastStatus)
	}
	if !repo.lastUpdatedAt.Equal(occurredAt) {
		t.Fatalf("expected updated_at %s, got %s", occurredAt, repo.lastUpdatedAt)
	}
	if result.FileID != fileID {
		t.Fatalf("expected result file id %s, got %s", fileID, result.FileID)
	}
}

func TestExecuteRejectsInvalidStatus(t *testing.T) {
	ctx := context.Background()
	uc := New(&processRepo{}, recordapp.Runtime{})

	_, err := uc.Execute(ctx, Command{
		FileID: uuid.New(),
		Status: sharedv1.FileStatus_FILE_STATUS_UPLOADING,
	})
	if err != recorderrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

type processRepo struct {
	attachment    domain.Attachment
	lastFileID    uuid.UUID
	lastStatus    sharedv1.FileStatus
	lastUpdatedAt time.Time
}

func (r *processRepo) CreateAttachment(context.Context, domain.AttachmentCreate, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *processRepo) AttachmentByID(context.Context, uuid.UUID, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *processRepo) ListAttachments(context.Context, domain.AttachmentListFilter) (domain.AttachmentList, error) {
	panic("unexpected call")
}

func (r *processRepo) UpdateAttachmentUploadStatusByFileID(_ context.Context, fileID uuid.UUID, status sharedv1.FileStatus, now time.Time) (domain.Attachment, error) {
	r.lastFileID = fileID
	r.lastStatus = status
	r.lastUpdatedAt = now
	r.attachment.FileID = fileID
	r.attachment.UploadStatus = status
	return r.attachment, nil
}

func (r *processRepo) DeleteAttachment(context.Context, uuid.UUID, uuid.UUID, time.Time) error {
	panic("unexpected call")
}

func (r *processRepo) PatientExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

func (r *processRepo) NoteExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

func (r *processRepo) NoteBelongsToPatient(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

var _ ports.Repository = (*processRepo)(nil)
