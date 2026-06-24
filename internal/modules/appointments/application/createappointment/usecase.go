package createappointment

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
	StartTime      time.Time
	EndTime        time.Time
	Reason         string
	PatientID      uuid.UUID
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.Appointment, error) {
	if cmd.PsychologistID == uuid.Nil {
		return domain.Appointment{}, recorderrors.ErrUnauthenticated
	}
	if cmd.PatientID == uuid.Nil || cmd.StartTime.IsZero() || cmd.EndTime.IsZero() || !cmd.EndTime.After(cmd.StartTime) || strings.TrimSpace(cmd.Reason) == "" {
		return domain.Appointment{}, recorderrors.ErrInvalidInput
	}
	exists, err := uc.repo.PatientExists(ctx, cmd.PsychologistID, cmd.PatientID)
	if err != nil {
		return domain.Appointment{}, err
	}
	if !exists {
		return domain.Appointment{}, recorderrors.ErrNotFound
	}
	return uc.repo.CreateAppointment(ctx, domain.AppointmentCreate{
		ID:             uc.newUUID(),
		PsychologistID: cmd.PsychologistID,
		PatientID:      cmd.PatientID,
		StartTime:      cmd.StartTime,
		EndTime:        cmd.EndTime,
		Reason:         strings.TrimSpace(cmd.Reason),
		Status:         sharedv1.AppointmentStatus_APPOINTMENT_STATUS_SCHEDULED,
		Now:            uc.now().UTC(),
	})
}
