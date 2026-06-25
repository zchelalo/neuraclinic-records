package listattachments

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	cfg   appshared.Config
	repo  ports.Repository
	files ports.FileManagementClient
}

func New(cfg appshared.Config, repo ports.Repository, files ports.FileManagementClient) *UseCase {
	return &UseCase{cfg: cfg, repo: repo, files: files}
}

type Command struct {
	PsychologistID uuid.UUID
	PatientID      uuid.UUID
	NoteID         *uuid.UUID
	Pagination     appshared.CursorPagination
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.AttachmentList, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.AttachmentList{}, recorderrors.ErrUnauthenticated
	}
	if cmd.PatientID == uuid.Nil {
		return domain.AttachmentList{}, recorderrors.ErrInvalidInput
	}
	exists, err := uc.repo.PatientExists(ctx, cmd.PsychologistID, cmd.PatientID)
	if err != nil {
		return domain.AttachmentList{}, err
	}
	if !exists {
		return domain.AttachmentList{}, recorderrors.ErrNotFound
	}
	pagination, err := appshared.ResolveCursorPagination(cmd.Pagination, uc.cfg)
	if err != nil {
		return domain.AttachmentList{}, recorderrors.ErrInvalidInput
	}
	result, err := uc.repo.ListAttachments(ctx, domain.AttachmentListFilter{
		PsychologistID: cmd.PsychologistID,
		PatientID:      cmd.PatientID,
		NoteID:         cmd.NoteID,
		Pagination:     pagination,
	})
	if err != nil {
		return domain.AttachmentList{}, err
	}
	for i := range result.Attachments {
		if result.Attachments[i].UploadStatus != sharedv1.FileStatus_FILE_STATUS_AVAILABLE {
			continue
		}
		if !isViewable(result.Attachments[i].MimeType) {
			continue
		}
		url, expiresAt, err := uc.files.GenerateDownloadURL(ctx, result.Attachments[i].FileID)
		if err != nil {
			if errors.Is(err, recorderrors.ErrFailedPrecondition) {
				continue
			}
			return domain.AttachmentList{}, err
		}
		result.Attachments[i].DownloadURL = &url
		result.Attachments[i].ExpiresAt = &expiresAt
	}
	return result, nil
}

func isViewable(mimeType string) bool {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	return strings.HasPrefix(mimeType, "image/") ||
		mimeType == "application/pdf" ||
		strings.HasPrefix(mimeType, "text/")
}
