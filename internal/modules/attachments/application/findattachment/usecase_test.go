package findattachment

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

func TestFindAttachmentSkipsDownloadURLWhenUploadIsNotAvailable(t *testing.T) {
	ctx := context.Background()
	attachmentID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	repo := &findRepo{
		attachment: domain.Attachment{
			ID:           attachmentID,
			FileID:       uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			MimeType:     "image/png",
			UploadStatus: sharedv1.FileStatus_FILE_STATUS_UPLOADING,
		},
	}
	files := &findFiles{}

	uc := New(repo, files)
	result, err := uc.Execute(ctx, Command{PsychologistID: uuid.New(), ID: attachmentID})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if files.calls != 0 {
		t.Fatalf("expected no download URL generation, got %d calls", files.calls)
	}
	if result.DownloadURL != nil {
		t.Fatal("expected nil download URL")
	}
}

func TestFindAttachmentGeneratesDownloadURLWhenAvailableAndViewable(t *testing.T) {
	ctx := context.Background()
	fileID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	repo := &findRepo{
		attachment: domain.Attachment{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			FileID:       fileID,
			MimeType:     "application/pdf",
			UploadStatus: sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
		},
	}
	files := &findFiles{
		url:       "http://download/file",
		expiresAt: time.Date(2026, 6, 25, 12, 0, 0, 0, time.UTC),
	}

	uc := New(repo, files)
	result, err := uc.Execute(ctx, Command{PsychologistID: uuid.New(), ID: repo.attachment.ID})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if files.calls != 1 {
		t.Fatalf("expected one download URL generation, got %d calls", files.calls)
	}
	if result.DownloadURL == nil || *result.DownloadURL != files.url {
		t.Fatalf("expected download url %q, got %#v", files.url, result.DownloadURL)
	}
}

type findRepo struct {
	attachment domain.Attachment
}

func (r *findRepo) CreateAttachment(context.Context, domain.AttachmentCreate, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *findRepo) AttachmentByID(context.Context, uuid.UUID, uuid.UUID) (domain.Attachment, error) {
	return r.attachment, nil
}

func (r *findRepo) ListAttachments(context.Context, domain.AttachmentListFilter) (domain.AttachmentList, error) {
	panic("unexpected call")
}

func (r *findRepo) UpdateAttachmentUploadStatusByFileID(context.Context, uuid.UUID, sharedv1.FileStatus, time.Time) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *findRepo) DeleteAttachment(context.Context, uuid.UUID, uuid.UUID, time.Time) error {
	panic("unexpected call")
}

func (r *findRepo) PatientExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

func (r *findRepo) NoteExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

func (r *findRepo) NoteBelongsToPatient(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

type findFiles struct {
	url       string
	expiresAt time.Time
	calls     int
}

func (f *findFiles) RequestUpload(context.Context, string, string, int64, bool, string) (uuid.UUID, string, time.Time, error) {
	panic("unexpected call")
}

func (f *findFiles) GenerateDownloadURL(context.Context, uuid.UUID) (string, time.Time, error) {
	f.calls++
	if f.url == "" {
		return "", time.Time{}, recorderrors.ErrFailedPrecondition
	}
	return f.url, f.expiresAt, nil
}

func (f *findFiles) Close() error {
	return nil
}

var _ ports.Repository = (*findRepo)(nil)
var _ ports.FileManagementClient = (*findFiles)(nil)
