package grpc

import (
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/application"
)

type FamiliogramService struct {
	recordv1.UnimplementedFamiliogramServiceServer
	app *application.Service
}

func NewFamiliogramService(app *application.Service) *FamiliogramService {
	return &FamiliogramService{app: app}
}
