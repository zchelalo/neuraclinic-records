package grpc

import (
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application"
)

type AttachmentService struct {
	recordv1.UnimplementedAttachmentServiceServer
	app *application.Service
}

func NewAttachmentService(app *application.Service) *AttachmentService {
	return &AttachmentService{app: app}
}
