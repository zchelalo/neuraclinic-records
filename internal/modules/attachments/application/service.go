package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application/createattachment"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application/deleteattachment"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application/findattachment"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application/listattachments"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application/processfilestatus"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
)

type Config = appshared.Config
type Runtime = appshared.Runtime

type Service struct {
	createAttachment  *createattachment.UseCase
	listAttachments   *listattachments.UseCase
	findAttachment    *findattachment.UseCase
	deleteAttachment  *deleteattachment.UseCase
	processFileStatus *processfilestatus.UseCase
}

func NewService(cfg Config, repo ports.Repository, files ports.FileManagementClient) *Service {
	return NewServiceWithRuntime(cfg, repo, files, appshared.DefaultRuntime())
}

func NewServiceWithRuntime(cfg Config, repo ports.Repository, files ports.FileManagementClient, runtime Runtime) *Service {
	runtime = runtime.Normalize()
	return &Service{
		createAttachment:  createattachment.New(repo, files, runtime),
		listAttachments:   listattachments.New(cfg, repo, files),
		findAttachment:    findattachment.New(repo, files),
		deleteAttachment:  deleteattachment.New(repo, runtime),
		processFileStatus: processfilestatus.New(repo, runtime),
	}
}

func DefaultRuntime() Runtime {
	return appshared.DefaultRuntime()
}

func (s *Service) CreateAttachment(ctx context.Context, cmd createattachment.Command) (domain.AttachmentCreateResult, error) {
	return s.createAttachment.Execute(ctx, cmd)
}

func (s *Service) ListAttachments(ctx context.Context, cmd listattachments.Command) (domain.AttachmentList, error) {
	return s.listAttachments.Execute(ctx, cmd)
}

func (s *Service) FindAttachment(ctx context.Context, psychologistID, id uuid.UUID) (domain.Attachment, error) {
	return s.findAttachment.Execute(ctx, findattachment.Command{PsychologistID: psychologistID, ID: id})
}

func (s *Service) DeleteAttachment(ctx context.Context, psychologistID, id uuid.UUID) error {
	return s.deleteAttachment.Execute(ctx, deleteattachment.Command{PsychologistID: psychologistID, ID: id})
}

func (s *Service) ProcessFileStatusChanged(ctx context.Context, cmd processfilestatus.Command) (domain.Attachment, error) {
	return s.processFileStatus.Execute(ctx, cmd)
}

type AttachmentCreateCommand = createattachment.Command
type AttachmentListCommand = listattachments.Command
