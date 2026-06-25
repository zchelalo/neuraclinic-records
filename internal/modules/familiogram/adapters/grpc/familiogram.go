package grpc

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	recordgrpc "github.com/zchelalo/neuraclinic-records/internal/shared/recordgrpc"
)

func (s *FamiliogramService) FindByPatientId(ctx context.Context, req *recordv1.FamiliogramServiceFindByPatientIdRequest) (*recordv1.FamiliogramServiceFindByPatientIdResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	patientID, err := recordgrpc.ParseID(req.GetPatientId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	familiogram, err := s.app.FindFamiliogram(ctx, psychologistID, patientID)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.FamiliogramServiceFindByPatientIdResponse{Familiogram: recordgrpc.FamiliogramToProto(familiogram)}, nil
}

func (s *FamiliogramService) Update(ctx context.Context, req *recordv1.FamiliogramServiceUpdateRequest) (*recordv1.FamiliogramServiceUpdateResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	familiogram, err := s.app.UpdateFamiliogram(ctx, psychologistID, id, req.GetData())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.FamiliogramServiceUpdateResponse{Familiogram: recordgrpc.FamiliogramToProto(familiogram)}, nil
}
