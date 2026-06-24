package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application/createappointment"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application/findappointment"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application/listappointments"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application/rescheduleappointment"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application/updateappointmentstatus"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
)

type Config = appshared.Config
type Runtime = appshared.Runtime

type Service struct {
	createAppointment       *createappointment.UseCase
	findAppointment         *findappointment.UseCase
	listAppointments        *listappointments.UseCase
	rescheduleAppointment   *rescheduleappointment.UseCase
	updateAppointmentStatus *updateappointmentstatus.UseCase
}

func NewService(cfg Config, repo ports.Repository) *Service {
	return NewServiceWithRuntime(cfg, repo, appshared.DefaultRuntime())
}

func NewServiceWithRuntime(cfg Config, repo ports.Repository, runtime Runtime) *Service {
	runtime = runtime.Normalize()
	return &Service{
		createAppointment:       createappointment.New(repo, runtime),
		findAppointment:         findappointment.New(repo),
		listAppointments:        listappointments.New(cfg, repo),
		rescheduleAppointment:   rescheduleappointment.New(repo, runtime),
		updateAppointmentStatus: updateappointmentstatus.New(repo, runtime),
	}
}

func DefaultRuntime() Runtime {
	return appshared.DefaultRuntime()
}

func (s *Service) CreateAppointment(ctx context.Context, cmd createappointment.Command) (domain.Appointment, error) {
	return s.createAppointment.Execute(ctx, cmd)
}

func (s *Service) FindAppointment(ctx context.Context, psychologistID, id uuid.UUID) (domain.Appointment, error) {
	return s.findAppointment.Execute(ctx, findappointment.Command{PsychologistID: psychologistID, ID: id})
}

func (s *Service) ListAppointments(ctx context.Context, cmd listappointments.Command) (domain.AppointmentList, error) {
	return s.listAppointments.Execute(ctx, cmd)
}

func (s *Service) RescheduleAppointment(ctx context.Context, cmd rescheduleappointment.Command) (domain.Appointment, error) {
	return s.rescheduleAppointment.Execute(ctx, cmd)
}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, cmd updateappointmentstatus.Command) (domain.Appointment, error) {
	return s.updateAppointmentStatus.Execute(ctx, cmd)
}

type AppointmentCreateCommand = createappointment.Command
type AppointmentListCommand = listappointments.Command
type AppointmentRescheduleCommand = rescheduleappointment.Command
type AppointmentStatusUpdateCommand = updateappointmentstatus.Command
