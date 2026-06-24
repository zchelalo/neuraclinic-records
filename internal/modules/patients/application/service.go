package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/application/createpatient"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/application/findpatient"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/application/listpatients"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/application/updatepatientaddress"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/application/updatepatientcontact"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/application/updatepatientidentification"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
)

type Config = appshared.Config
type Runtime = appshared.Runtime

type Service struct {
	createPatient               *createpatient.UseCase
	listPatients                *listpatients.UseCase
	findPatient                 *findpatient.UseCase
	updatePatientIdentification *updatepatientidentification.UseCase
	updatePatientContact        *updatepatientcontact.UseCase
	updatePatientAddress        *updatepatientaddress.UseCase
}

func NewService(cfg Config, repo ports.Repository) *Service {
	return NewServiceWithRuntime(cfg, repo, appshared.DefaultRuntime())
}

func NewServiceWithRuntime(cfg Config, repo ports.Repository, runtime Runtime) *Service {
	runtime = runtime.Normalize()
	return &Service{
		createPatient:               createpatient.New(repo, runtime),
		listPatients:                listpatients.New(cfg, repo),
		findPatient:                 findpatient.New(repo),
		updatePatientIdentification: updatepatientidentification.New(repo, runtime),
		updatePatientContact:        updatepatientcontact.New(repo, runtime),
		updatePatientAddress:        updatepatientaddress.New(repo, runtime),
	}
}

func DefaultRuntime() Runtime {
	return appshared.DefaultRuntime()
}

func (s *Service) CreatePatient(ctx context.Context, cmd createpatient.Command) (domain.PatientSummary, error) {
	return s.createPatient.Execute(ctx, cmd)
}

func (s *Service) ListPatients(ctx context.Context, cmd listpatients.Command) (domain.PatientList, error) {
	return s.listPatients.Execute(ctx, cmd)
}

func (s *Service) FindPatient(ctx context.Context, psychologistID, id uuid.UUID) (domain.Patient, error) {
	return s.findPatient.Execute(ctx, findpatient.Command{PsychologistID: psychologistID, ID: id})
}

func (s *Service) UpdatePatientIdentification(ctx context.Context, cmd updatepatientidentification.Command) (domain.Patient, error) {
	return s.updatePatientIdentification.Execute(ctx, cmd)
}

func (s *Service) UpdatePatientContact(ctx context.Context, cmd updatepatientcontact.Command) (domain.Patient, error) {
	return s.updatePatientContact.Execute(ctx, cmd)
}

func (s *Service) UpdatePatientAddress(ctx context.Context, cmd updatepatientaddress.Command) (domain.Patient, error) {
	return s.updatePatientAddress.Execute(ctx, cmd)
}

type PatientCreateCommand = createpatient.Command
type PatientListCommand = listpatients.Command
type PatientIdentificationUpdateCommand = updatepatientidentification.Command
type PatientContactUpdateCommand = updatepatientcontact.Command
type PatientAddressUpdateCommand = updatepatientaddress.Command
