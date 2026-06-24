package updateappointmentstatus

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	repo ports.Repository
	now  func() time.Time
}

func New(repo ports.Repository, runtime appshared.Runtime) *UseCase {
	runtime = runtime.Normalize()
	return &UseCase{repo: repo, now: runtime.Now}
}

type Command struct {
	PsychologistID    uuid.UUID
	AppointmentID     uuid.UUID
	Status            sharedv1.AppointmentStatus
	CancelledByUserID *uuid.UUID
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Appointment, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Appointment{}, recorderrors.ErrUnauthenticated
	}
	if cmd.AppointmentID == uuid.Nil {
		return domain.Appointment{}, recorderrors.ErrInvalidInput
	}
	switch cmd.Status {
	case sharedv1.AppointmentStatus_APPOINTMENT_STATUS_COMPLETED,
		sharedv1.AppointmentStatus_APPOINTMENT_STATUS_CANCELLED,
		sharedv1.AppointmentStatus_APPOINTMENT_STATUS_NO_SHOW:
	default:
		return domain.Appointment{}, recorderrors.ErrInvalidInput
	}
	if cmd.Status == sharedv1.AppointmentStatus_APPOINTMENT_STATUS_CANCELLED && (cmd.CancelledByUserID == nil || *cmd.CancelledByUserID == uuid.Nil) {
		return domain.Appointment{}, recorderrors.ErrInvalidInput
	}
	return uc.repo.UpdateAppointmentStatus(ctx, cmd.PsychologistID, cmd.AppointmentID, cmd.Status, cmd.CancelledByUserID, uc.now().UTC())
}
