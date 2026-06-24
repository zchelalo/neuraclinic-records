package findattachment

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	repo  ports.Repository
	files ports.FileManagementClient
}

func New(repo ports.Repository, files ports.FileManagementClient) *UseCase {
	return &UseCase{repo: repo, files: files}
}

type Command struct {
	PsychologistID uuid.UUID
	ID             uuid.UUID
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Attachment, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Attachment{}, recorderrors.ErrUnauthenticated
	}
	if cmd.ID == uuid.Nil {
		return domain.Attachment{}, recorderrors.ErrInvalidInput
	}
	attachment, err := uc.repo.AttachmentByID(ctx, cmd.PsychologistID, cmd.ID)
	if err != nil {
		return domain.Attachment{}, err
	}
	if isViewable(attachment.MimeType) {
		url, expiresAt, err := uc.files.GenerateDownloadURL(ctx, attachment.FileID)
		if err != nil {
			return domain.Attachment{}, err
		}
		attachment.DownloadURL = &url
		attachment.ExpiresAt = &expiresAt
	}
	return attachment, nil
}

func isViewable(mimeType string) bool {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	return strings.HasPrefix(mimeType, "image/") ||
		mimeType == "application/pdf" ||
		strings.HasPrefix(mimeType, "text/")
}
