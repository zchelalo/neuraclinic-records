package rabbitmq

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	filemanagementv1 "github.com/zchelalo/neuraclinic-records/gen/go/file_management/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	recordapp "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandlerAcksMissingAttachment(t *testing.T) {
	service := application.NewServiceWithRuntime(recordapp.Config{}, &handlerRepo{err: recorderrors.ErrNotFound}, &handlerFiles{}, recordapp.Runtime{})
	handler := NewHandler(service)

	body, err := protojson.Marshal(&filemanagementv1.FileStatusChangedEvent{
		EventId:       uuid.NewString(),
		FileId:        uuid.NewString(),
		ServiceOrigin: "record",
		Status:        sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
		OccurredAt:    timestamppb.New(time.Date(2026, 6, 25, 18, 0, 0, 0, time.UTC)),
	})
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	if err := handler(context.Background(), "file.record.status_changed.v1", body); err != nil {
		t.Fatalf("expected nil error for missing attachment, got %v", err)
	}
}

func TestHandlerRejectsInvalidEvent(t *testing.T) {
	service := application.NewServiceWithRuntime(recordapp.Config{}, &handlerRepo{}, &handlerFiles{}, recordapp.Runtime{})
	handler := NewHandler(service)

	err := handler(context.Background(), "file.record.status_changed.v1", []byte(`{"event_id":"bad"}`))
	if !errors.Is(err, ports.ErrInvalidEvent) {
		t.Fatalf("expected invalid event error, got %v", err)
	}
}

type handlerRepo struct {
	err error
}

func (r *handlerRepo) CreateAttachment(context.Context, domain.AttachmentCreate, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *handlerRepo) AttachmentByID(context.Context, uuid.UUID, uuid.UUID) (domain.Attachment, error) {
	panic("unexpected call")
}

func (r *handlerRepo) ListAttachments(context.Context, domain.AttachmentListFilter) (domain.AttachmentList, error) {
	panic("unexpected call")
}

func (r *handlerRepo) UpdateAttachmentUploadStatusByFileID(context.Context, uuid.UUID, sharedv1.FileStatus, time.Time) (domain.Attachment, error) {
	if r.err != nil {
		return domain.Attachment{}, r.err
	}
	return domain.Attachment{}, nil
}

func (r *handlerRepo) DeleteAttachment(context.Context, uuid.UUID, uuid.UUID, time.Time) error {
	panic("unexpected call")
}

func (r *handlerRepo) PatientExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

func (r *handlerRepo) NoteBelongsToPatient(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (bool, error) {
	panic("unexpected call")
}

type handlerFiles struct{}

func (f *handlerFiles) RequestUpload(context.Context, string, string, int64, bool, string) (uuid.UUID, string, time.Time, error) {
	panic("unexpected call")
}

func (f *handlerFiles) GenerateDownloadURL(context.Context, uuid.UUID) (string, time.Time, error) {
	panic("unexpected call")
}

func (f *handlerFiles) Close() error {
	return nil
}

var _ ports.Repository = (*handlerRepo)(nil)
var _ ports.FileManagementClient = (*handlerFiles)(nil)
