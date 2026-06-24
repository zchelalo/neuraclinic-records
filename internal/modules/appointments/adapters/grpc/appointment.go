package grpc

import (
	"context"
	"time"

	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/appointments/application"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	recordgrpc "github.com/zchelalo/neuraclinic-records/internal/shared/recordgrpc"
)

func (s *AppointmentService) Create(ctx context.Context, req *recordv1.AppointmentServiceCreateRequest) (*recordv1.AppointmentServiceCreateResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	patientID, err := recordgrpc.ParseID(req.GetPatientId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	if req.GetStartTime() == nil || req.GetEndTime() == nil {
		return nil, recordgrpc.MapError(recorderrors.ErrInvalidInput)
	}
	appointment, err := s.app.CreateAppointment(ctx, application.AppointmentCreateCommand{
		PsychologistID: psychologistID,
		StartTime:      req.GetStartTime().AsTime(),
		EndTime:        req.GetEndTime().AsTime(),
		Reason:         req.GetReason(),
		PatientID:      patientID,
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.AppointmentServiceCreateResponse{Appointment: recordgrpc.AppointmentToProto(appointment)}, nil
}

func (s *AppointmentService) FindById(ctx context.Context, req *recordv1.AppointmentServiceFindByIdRequest) (*recordv1.AppointmentServiceFindByIdResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetAppointmentId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	appointment, err := s.app.FindAppointment(ctx, psychologistID, id)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.AppointmentServiceFindByIdResponse{Appointment: recordgrpc.AppointmentToProto(appointment)}, nil
}

func (s *AppointmentService) List(ctx context.Context, req *recordv1.AppointmentServiceListRequest) (*recordv1.AppointmentServiceListResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	patientID, err := recordgrpc.ParseOptionalID(req.PatientId)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	var startDate, endDate *time.Time
	if req.GetDateRange() != nil {
		if req.GetDateRange().GetStartDate() != nil {
			value := req.GetDateRange().GetStartDate().AsTime()
			startDate = &value
		}
		if req.GetDateRange().GetEndDate() != nil {
			value := req.GetDateRange().GetEndDate().AsTime()
			endDate = &value
		}
	}
	result, err := s.app.ListAppointments(ctx, application.AppointmentListCommand{
		PsychologistID: psychologistID,
		Pagination:     recordgrpc.CursorPaginationFromProto(req.GetPagination()),
		PatientID:      patientID,
		StartDate:      startDate,
		EndDate:        endDate,
		Statuses:       req.GetStatuses(),
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	appointments := make([]*recordv1.Appointment, 0, len(result.Appointments))
	for _, appointment := range result.Appointments {
		appointments = append(appointments, recordgrpc.AppointmentToProto(appointment))
	}
	return &recordv1.AppointmentServiceListResponse{
		Appointments: appointments,
		Meta:         recordgrpc.CursorMetaToProto(result.Meta),
	}, nil
}

func (s *AppointmentService) Reschedule(ctx context.Context, req *recordv1.AppointmentServiceRescheduleRequest) (*recordv1.AppointmentServiceRescheduleResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetAppointmentId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	if req.GetNewStartTime() == nil || req.GetNewEndTime() == nil {
		return nil, recordgrpc.MapError(recorderrors.ErrInvalidInput)
	}
	appointment, err := s.app.RescheduleAppointment(ctx, application.AppointmentRescheduleCommand{
		PsychologistID: psychologistID,
		AppointmentID:  id,
		NewStartTime:   req.GetNewStartTime().AsTime(),
		NewEndTime:     req.GetNewEndTime().AsTime(),
		Reason:         req.GetReason(),
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.AppointmentServiceRescheduleResponse{Appointment: recordgrpc.AppointmentToProto(appointment)}, nil
}

func (s *AppointmentService) UpdateStatus(ctx context.Context, req *recordv1.AppointmentServiceUpdateStatusRequest) (*recordv1.AppointmentServiceUpdateStatusResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	id, err := recordgrpc.ParseID(req.GetAppointmentId())
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	cancelledByUserID, err := recordgrpc.ParseOptionalID(req.CancelledByUserId)
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	appointment, err := s.app.UpdateAppointmentStatus(ctx, application.AppointmentStatusUpdateCommand{
		PsychologistID:    psychologistID,
		AppointmentID:     id,
		Status:            req.GetNewStatus(),
		CancelledByUserID: cancelledByUserID,
	})
	if err != nil {
		return nil, recordgrpc.MapError(err)
	}
	return &recordv1.AppointmentServiceUpdateStatusResponse{Appointment: recordgrpc.AppointmentToProto(appointment)}, nil
}
