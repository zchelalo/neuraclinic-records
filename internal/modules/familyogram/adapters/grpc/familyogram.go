package grpc

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	recordgrpc "github.com/zchelalo/neuraclinic-records/internal/shared/recordgrpc"
)

func (s *FamilyogramService) FindByPatientId(ctx context.Context, req *recordv1.FamilyogramServiceFindByPatientIdRequest) (*recordv1.FamilyogramServiceFindByPatientIdResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	patientID, err := recordgrpc.ParseID(req.GetPatientId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	familyogram, err := s.app.FindFamilyogram(ctx, psychologistID, patientID)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.FamilyogramServiceFindByPatientIdResponse{Familyogram: recordgrpc.FamilyogramToProto(familyogram)}, nil
}

func (s *FamilyogramService) Update(ctx context.Context, req *recordv1.FamilyogramServiceUpdateRequest) (*recordv1.FamilyogramServiceUpdateResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	familyogram, err := s.app.UpdateFamilyogram(ctx, psychologistID, id, req.GetData())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.FamilyogramServiceUpdateResponse{Familyogram: recordgrpc.FamilyogramToProto(familyogram)}, nil
}
