package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/application/findfamilyogram"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/application/updatefamilyogram"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/domain"
	"github.com/zchelalo/neuraclinic-records/internal/modules/familyogram/ports"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	"google.golang.org/protobuf/types/known/structpb"
)

type Runtime = appshared.Runtime

type Service struct {
	findFamilyogram   *findfamilyogram.UseCase
	updateFamilyogram *updatefamilyogram.UseCase
}

func NewService(repo ports.Repository) *Service {
	return NewServiceWithRuntime(repo, appshared.DefaultRuntime())
}

func NewServiceWithRuntime(repo ports.Repository, runtime Runtime) *Service {
	runtime = runtime.Normalize()
	return &Service{
		findFamilyogram:   findfamilyogram.New(repo),
		updateFamilyogram: updatefamilyogram.New(repo, runtime),
	}
}

func DefaultRuntime() Runtime {
	return appshared.DefaultRuntime()
}

func (s *Service) FindFamilyogram(ctx context.Context, psychologistID, patientID uuid.UUID) (domain.Familyogram, error) {
	return s.findFamilyogram.Execute(ctx, findfamilyogram.Command{PsychologistID: psychologistID, PatientID: patientID})
}

func (s *Service) UpdateFamilyogram(ctx context.Context, psychologistID, id uuid.UUID, data *structpb.Struct) (domain.Familyogram, error) {
	return s.updateFamilyogram.Execute(ctx, updatefamilyogram.Command{PsychologistID: psychologistID, ID: id, Data: data})
}
