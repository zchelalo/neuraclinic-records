package createattachment

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

func TestExecuteReturnsNotFoundWhenNoteDoesNotExist(t *testing.T) {
	ctx := context.Background()
	repo := &createRepo{
		patientExists: true,
		noteExists:    false,
	}

	uc := New(repo, &createFiles{}, appshared.Runtime{})
	_, err := uc.Execute(ctx, Command{
		PsychologistID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		PatientID:      uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		NoteID:         ptrUUID(uuid.MustParse("33333333-3333-3333-3333-333333333333")),
		OriginalName:   "evidence.pdf",
		MimeType:       "application/pdf",
		SizeBytes:      1024,
	})
	if err != recorderrors.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if repo.noteBelongsCalls != 0 {
		t.Fatalf("expected note ownership check to be skipped, got %d calls", repo.noteBelongsCalls)
	}
}

func TestExecuteReturnsInvalidInputWhenNoteBelongsToAnotherPatient(t *testing.T) {
	ctx := context.Background()
	repo := &createRepo{
		patientExists:      true,
		noteExists:         true,
		noteBelongsToPatient: false,
	}

	uc := New(repo, &createFiles{}, appshared.Runtime{})
	_, err := uc.Execute(ctx, Command{
		PsychologistID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		PatientID:      uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		NoteID:         ptrUUID(uuid.MustParse("33333333-3333-3333-3333-333333333333")),
		OriginalName:   "evidence.pdf",
		MimeType:       "application/pdf",
		SizeBytes:      1024,
	})
	if err != recorderrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

type createRepo struct {
	patientExists        bool
	noteExists           bool
	noteBelongsToPatient bool
	noteBelongsCalls     int
}

func (r *createRepo) CreateAttachment(context.Context, domain.AttachmentCreate, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *createRepo) AttachmentByID(context.Context, uuid.UUID, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *createRepo) ListAttachments(context.Context, domain.AttachmentListFilter) (domain.AttachmentList, error) {
	panic("unexpected call")
}

func (r *createRepo) UpdateAttachmentUploadStatusByFileID(context.Context, uuid.UUID, sharedv1.FileStatus, time.Time) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *createRepo) DeleteAttachment(context.Context, uuid.UUID, uuid.UUID, time.Time) error {
	panic("unexpected call")
}

func (r *createRepo) PatientExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return r.patientExists, nil
}

func (r *createRepo) NoteExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return r.noteExists, nil
}

func (r *createRepo) NoteBelongsToPatient(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (bool, error) {
	r.noteBelongsCalls++
	return r.noteBelongsToPatient, nil
}

type createFiles struct{}

func (f *createFiles) RequestUpload(context.Context, string, string, int64, bool, string) (uuid.UUID, string, time.Time, error) {
	panic("unexpected call")
}

func (f *createFiles) GenerateDownloadURL(context.Context, uuid.UUID) (string, time.Time, error) {
	panic("unexpected call")
}

func (f *createFiles) Close() error {
	return nil
}

func ptrUUID(id uuid.UUID) *uuid.UUID {
	return &id
}

var _ ports.Repository = (*createRepo)(nil)
var _ ports.FileManagementClient = (*createFiles)(nil)
