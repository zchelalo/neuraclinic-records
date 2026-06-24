package grpc

import (
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/application"
)

type FamilyogramService struct {
	recordv1.UnimplementedFamilyogramServiceServer
	app *application.Service
}

func NewFamilyogramService(app *application.Service) *FamilyogramService {
	return &FamilyogramService{app: app}
}
