package listattachments

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

func TestListAttachmentsOnlyGeneratesDownloadURLForAvailableViewableFiles(t *testing.T) {
	ctx := context.Background()
	psychologistID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	patientID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	availableID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	uploadingID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	archiveID := uuid.MustParse("55555555-5555-5555-5555-555555555555")
	repo := &listRepo{
		patientExists: true,
		list: domain.AttachmentList{
			Attachments: []domain.Attachment{
				{ID: uuid.New(), FileID: availableID, MimeType: "image/png", UploadStatus: sharedv1.FileStatus_FILE_STATUS_AVAILABLE},
				{ID: uuid.New(), FileID: uploadingID, MimeType: "image/png", UploadStatus: sharedv1.FileStatus_FILE_STATUS_UPLOADING},
				{ID: uuid.New(), FileID: archiveID, MimeType: "application/zip", UploadStatus: sharedv1.FileStatus_FILE_STATUS_AVAILABLE},
			},
		},
	}
	files := &listFiles{
		urls: map[uuid.UUID]string{
			availableID: "http://download/available",
		},
		expiresAt: time.Date(2026, 6, 25, 10, 0, 0, 0, time.UTC),
	}

	uc := New(appshared.Config{PaginationLimitDefault: 10, PaginationLimitMax: 100}, repo, files)
	result, err := uc.Execute(ctx, Command{
		PsychologistID: psychologistID,
		PatientID:      patientID,
		Pagination:     appshared.CursorPagination{Limit: 10},
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if len(result.Attachments) != 2 {
		t.Fatalf("expected only available attachments in response, got %d", len(result.Attachments))
	}
	if files.calls != 1 {
		t.Fatalf("expected one download URL generation, got %d", files.calls)
	}
	if result.Attachments[0].UploadStatus != sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
		t.Fatalf("expected first attachment to be available, got %s", result.Attachments[0].UploadStatus)
	}
	if result.Attachments[0].DownloadURL == nil || *result.Attachments[0].DownloadURL != "http://download/available" {
		t.Fatalf("expected available attachment download url, got %#v", result.Attachments[0].DownloadURL)
	}
	if result.Attachments[1].UploadStatus != sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
		t.Fatalf("expected second attachment to be available, got %s", result.Attachments[1].UploadStatus)
	}
	if result.Attachments[1].DownloadURL != nil {
		t.Fatal("expected non-viewable attachment to have no download url")
	}
}

type listRepo struct {
	patientExists bool
	list          domain.AttachmentList
}

func (r *listRepo) CreateAttachment(context.Context, domain.AttachmentCreate, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *listRepo) AttachmentByID(context.Context, uuid.UUID, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *listRepo) ListAttachments(context.Context, domain.AttachmentListFilter) (domain.AttachmentList, error) {
	return r.list, nil
}

func (r *listRepo) UpdateAttachmentUploadStatusByFileID(context.Context, uuid.UUID, sharedv1.FileStatus, time.Time) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *listRepo) DeleteAttachment(context.Context, uuid.UUID, uuid.UUID, time.Time) error {
	panic("unexpected call")
}

func (r *listRepo) PatientExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return r.patientExists, nil
}

func (r *listRepo) NoteExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

func (r *listRepo) NoteBelongsToPatient(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

type listFiles struct {
	urls      map[uuid.UUID]string
	expiresAt time.Time
	calls     int
}

func (f *listFiles) RequestUpload(context.Context, string, string, int64, bool, string) (uuid.UUID, string, time.Time, error) {
	panic("unexpected call")
}

func (f *listFiles) GenerateDownloadURL(_ context.Context, id uuid.UUID) (string, time.Time, error) {
	f.calls++
	url, ok := f.urls[id]
	if !ok {
		return "", time.Time{}, recorderrors.ErrFailedPrecondition
	}
	return url, f.expiresAt, nil
}

func (f *listFiles) Close() error {
	return nil
}

var _ ports.Repository = (*listRepo)(nil)
var _ ports.FileManagementClient = (*listFiles)(nil)
