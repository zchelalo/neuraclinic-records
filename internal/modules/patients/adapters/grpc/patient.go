package grpc

import (
	"context"
	"time"

	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/patients/application"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	recordgrpc "github.com/zchelalo/neuraclinic-records/internal/shared/recordgrpc"
)

func (s *PatientService) Create(ctx context.Context, req *recordv1.PatientServiceCreateRequest) (*recordv1.PatientServiceCreateResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	birthDate, err := recordgrpc.TimeFromProtoDate(req.GetBirthDate())
	if err != nil {
		return nil, recordgrpc.MapError(recorderrors.ErrInvalidInput)
	}
	patient, err := s.app.CreatePatient(ctx, application.PatientCreateCommand{
		PsychologistID: psychologistID,
		FirstName:      req.GetFirstName(),
		MiddleName:     req.MiddleName,
		FirstLastName:  req.GetFirstLastName(),
		SecondLastName: req.SecondLastName,
		BirthDate:      birthDate,
		BirthCountry:   req.GetBirthCountry(),
		BirthProvince:  req.GetBirthProvince(),
		BirthCity:      req.GetBirthCity(),
		Sex:            req.GetSex(),
		MaritalStatus:  req.GetMaritalStatus(),
		Occupation:     req.Occupation,
		Religion:       req.Religion,
		Phone:          req.GetPhone(),
		Email:          req.GetEmail(),
		Country:        req.GetCountry(),
		Province:       req.GetProvince(),
		City:           req.GetCity(),
		PostalCode:     req.GetPostalCode(),
		Neighborhood:   req.GetNeighborhood(),
		Street:         req.GetStreet(),
		StreetNumber:   req.GetStreetNumber(),
		UnitNumber:     req.UnitNumber,
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.PatientServiceCreateResponse{Patient: recordgrpc.PatientSummaryToProto(patient)}, nil
}

func (s *PatientService) List(ctx context.Context, req *recordv1.PatientServiceListRequest) (*recordv1.PatientServiceListResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	result, err := s.app.ListPatients(ctx, application.PatientListCommand{
		PsychologistID:          psychologistID,
		Pagination:              recordgrpc.CursorPaginationFromProto(req.GetPagination()),
		WithPendingAppointments: req.GetWithPendingAppointments(),
		WithNoAppointments:      req.GetWithNoAppointments(),
		EverHadAppointments:     req.GetEverHadAppointments(),
		SearchQuery:             req.GetSearchQuery(),
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	patients := make([]*recordv1.PatientSummary, 0, len(result.Patients))
	for _, patient := range result.Patients {
		patients = append(patients, recordgrpc.PatientSummaryToProto(patient))
	}
	return &recordv1.PatientServiceListResponse{
		Patients: patients,
		Meta:     recordgrpc.CursorMetaToProto(result.Meta),
	}, nil
}

func (s *PatientService) FindById(ctx context.Context, req *recordv1.PatientServiceFindByIdRequest) (*recordv1.PatientServiceFindByIdResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	patient, err := s.app.FindPatient(ctx, psychologistID, id)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.PatientServiceFindByIdResponse{Patient: recordgrpc.PatientToProto(patient)}, nil
}

func (s *PatientService) UpdateIdentificationData(ctx context.Context, req *recordv1.PatientServiceUpdateIdentificationDataRequest) (*recordv1.PatientServiceUpdateIdentificationDataResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	var birthDate *time.Time
	if req.BirthDate != nil {
		parsed, err := recordgrpc.TimeFromProtoDate(req.BirthDate)
		if err != nil {
			return nil, recordgrpc.MapError(recorderrors.ErrInvalidInput)
		}
		birthDate = &parsed
	}
	patient, err := s.app.UpdatePatientIdentification(ctx, application.PatientIdentificationUpdateCommand{
		PsychologistID: psychologistID,
		ID:             id,
		FirstName:      req.FirstName,
		MiddleName:     req.MiddleName,
		FirstLastName:  req.FirstLastName,
		SecondLastName: req.SecondLastName,
		BirthDate:      birthDate,
		Sex:            req.Sex,
		BirthCountry:   req.BirthCountry,
		BirthProvince:  req.BirthProvince,
		BirthCity:      req.BirthCity,
		Occupation:     req.Occupation,
		MaritalStatus:  req.MaritalStatus,
		Religion:       req.Religion,
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.PatientServiceUpdateIdentificationDataResponse{Patient: recordgrpc.PatientToProto(patient)}, nil
}

func (s *PatientService) UpdateContactDetails(ctx context.Context, req *recordv1.PatientServiceUpdateContactDetailsRequest) (*recordv1.PatientServiceUpdateContactDetailsResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	patient, err := s.app.UpdatePatientContact(ctx, application.PatientContactUpdateCommand{
		PsychologistID: psychologistID,
		ID:             id,
		Phone:          req.Phone,
		Email:          req.Email,
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.PatientServiceUpdateContactDetailsResponse{Patient: recordgrpc.PatientToProto(patient)}, nil
}

func (s *PatientService) UpdateAddress(ctx context.Context, req *recordv1.PatientServiceUpdateAddressRequest) (*recordv1.PatientServiceUpdateAddressResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	patient, err := s.app.UpdatePatientAddress(ctx, application.PatientAddressUpdateCommand{
		PsychologistID: psychologistID,
		ID:             id,
		Country:        req.Country,
		Province:       req.Province,
		City:           req.City,
		PostalCode:     req.PostalCode,
		Neighborhood:   req.Neighborhood,
		Street:         req.Street,
		StreetNumber:   req.StreetNumber,
		UnitNumber:     req.UnitNumber,
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.PatientServiceUpdateAddressResponse{Patient: recordgrpc.PatientToProto(patient)}, nil
}
