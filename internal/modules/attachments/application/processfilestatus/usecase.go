package processfilestatus

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	repo ports.Repository
	now  func() time.Time
}

func New(repo ports.Repository, runtime appshared.Runtime) *UseCase {
	runtime = runtime.Normalize()
	return &UseCase{repo: repo, now: runtime.Now}
}

type Command struct {
	FileID     uuid.UUID
	Status     sharedv1.FileStatus
	OccurredAt *time.Time
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Attachment, error) {
	if cmd.FileID == uuid.Nil {
		return domain.Attachment{}, recorderrors.ErrInvalidInput
	}
	if cmd.Status != sharedv1.FileStatus_FILE_STATUS_AVAILABLE && cmd.Status != sharedv1.FileStatus_FILE_STATUS_ERROR {
		return domain.Attachment{}, recorderrors.ErrInvalidInput
	}

	updatedAt := uc.now().UTC()
	if cmd.OccurredAt != nil && !cmd.OccurredAt.IsZero() {
		updatedAt = cmd.OccurredAt.UTC()
	}

	return uc.repo.UpdateAttachmentUploadStatusByFileID(ctx, cmd.FileID, cmd.Status, updatedAt)
}
