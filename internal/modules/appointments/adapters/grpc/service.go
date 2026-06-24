package grpc

import (
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application"
)

type AppointmentService struct {
	recordv1.UnimplementedAppointmentServiceServer
	app *application.Service
}

func NewAppointmentService(app *application.Service) *AppointmentService {
	return &AppointmentService{app: app}
}
