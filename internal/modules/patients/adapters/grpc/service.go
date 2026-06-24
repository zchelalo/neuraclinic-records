package grpc

import (
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/application"
)

type PatientService struct {
	recordv1.UnimplementedPatientServiceServer
	app *application.Service
}

func NewPatientService(app *application.Service) *PatientService {
	return &PatientService{app: app}
}
