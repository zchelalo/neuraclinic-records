package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/application/findfamiliogram"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/application/updatefamiliogram"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	"google.golang.org/protobuf/types/known/structpb"
)

type Runtime = appshared.Runtime

type Service struct {
	findFamiliogram   *findfamiliogram.UseCase
	updateFamiliogram *updatefamiliogram.UseCase
}

func NewService(repo ports.Repository) *Service {
	return NewServiceWithRuntime(repo, appshared.DefaultRuntime())
}

func NewServiceWithRuntime(repo ports.Repository, runtime Runtime) *Service {
	runtime = runtime.Normalize()
	return &Service{
		findFamiliogram:   findfamiliogram.New(repo),
		updateFamiliogram: updatefamiliogram.New(repo, runtime),
	}
}

func DefaultRuntime() Runtime {
	return appshared.DefaultRuntime()
}

func (s *Service) FindFamiliogram(ctx context.Context, psychologistID, patientID uuid.UUID) (domain.Familiogram, error) {
	return s.findFamiliogram.Execute(ctx, findfamiliogram.Command{PsychologistID: psychologistID, PatientID: patientID})
}

func (s *Service) UpdateFamiliogram(ctx context.Context, psychologistID, id uuid.UUID, data *structpb.Struct) (domain.Familiogram, error) {
	return s.updateFamiliogram.Execute(ctx, updatefamiliogram.Command{PsychologistID: psychologistID, ID: id, Data: data})
}
