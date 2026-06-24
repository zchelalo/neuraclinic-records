package rescheduleappointment

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
)

type UseCase struct {
	repo    ports.Repository
	now     func() time.Time
	newUUID func() uuid.UUID
}

func New(repo ports.Repository, runtime appshared.Runtime) *UseCase {
	runtime = runtime.Normalize()
	return &UseCase{repo: repo, now: runtime.Now, newUUID: runtime.NewUUID}
}

type Command struct {
	PsychologistID uuid.UUID
	AppointmentID  uuid.UUID
	NewStartTime   time.Time
	NewEndTime     time.Time
	Reason         string
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Appointment, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Appointment{}, recorderrors.ErrUnauthenticated
	}
	if cmd.AppointmentID == uuid.Nil || cmd.NewStartTime.IsZero() || cmd.NewEndTime.IsZero() || !cmd.NewEndTime.After(cmd.NewStartTime) || strings.TrimSpace(cmd.Reason) == "" {
		return domain.Appointment{}, recorderrors.ErrInvalidInput
	}
	original, err := uc.repo.AppointmentByID(ctx, cmd.PsychologistID, cmd.AppointmentID)
	if err != nil {
		return domain.Appointment{}, err
	}
	if original.Status != sharedv1.AppointmentStatus_APPOINTMENT_STATUS_SCHEDULED {
		return domain.Appointment{}, recorderrors.ErrFailedPrecondition
	}
	now := uc.now().UTC()
	return uc.repo.RescheduleAppointment(ctx, cmd.PsychologistID, cmd.AppointmentID, domain.AppointmentCreate{
		ID:                           uc.newUUID(),
		PsychologistID:               cmd.PsychologistID,
		PatientID:                    original.PatientID,
		StartTime:                    cmd.NewStartTime,
		EndTime:                      cmd.NewEndTime,
		Reason:                       strings.TrimSpace(cmd.Reason),
		Status:                       sharedv1.AppointmentStatus_APPOINTMENT_STATUS_SCHEDULED,
		RescheduledFromAppointmentID: &cmd.AppointmentID,
		Now:                          now,
	}, now)
}
