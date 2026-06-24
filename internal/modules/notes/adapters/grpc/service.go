package grpc

import (
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/application"
)

type NoteService struct {
	recordv1.UnimplementedNoteServiceServer
	app *application.Service
}

func NewNoteService(app *application.Service) *NoteService {
	return &NoteService{app: app}
}
