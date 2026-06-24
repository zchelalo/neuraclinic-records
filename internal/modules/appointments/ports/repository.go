package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/domain"
)

type Repository interface {
	CreateAppointment(ctx context.Context, appointment domain.AppointmentCreate) (domain.Appointment, error)
	AppointmentByID(ctx context.Context, psychologistID, id uuid.UUID) (domain.Appointment, error)
	ListAppointments(ctx context.Context, filter domain.AppointmentListFilter) (domain.AppointmentList, error)
	RescheduleAppointment(ctx context.Context, psychologistID, originalID uuid.UUID, appointment domain.AppointmentCreate, now time.Time) (domain.Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, psychologistID, id uuid.UUID, status sharedv1.AppointmentStatus, cancelledByUserID *uuid.UUID, now time.Time) (domain.Appointment, error)
	AppointmentBelongsToPatient(ctx context.Context, psychologistID, appointmentID, patientID uuid.UUID) (bool, error)
	PatientExists(ctx context.Context, psychologistID, id uuid.UUID) (bool, error)
}
