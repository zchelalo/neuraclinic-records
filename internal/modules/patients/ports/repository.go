package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/domain"
)

type Repository interface {
	CreatePatient(ctx context.Context, patient domain.PatientCreate) (domain.PatientSummary, error)
	ListPatients(ctx context.Context, filter domain.PatientListFilter) (domain.PatientList, error)
	PatientByID(ctx context.Context, psychologistID, id uuid.UUID) (domain.Patient, error)
	UpdatePatientIdentification(ctx context.Context, update domain.PatientIdentificationUpdate) (domain.Patient, error)
	UpdatePatientContact(ctx context.Context, update domain.PatientContactUpdate) (domain.Patient, error)
	UpdatePatientAddress(ctx context.Context, update domain.AddressUpdate) (domain.Patient, error)
	PatientExists(ctx context.Context, psychologistID, id uuid.UUID) (bool, error)
}
