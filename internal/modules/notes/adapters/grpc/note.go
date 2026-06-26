package grpc

import (
	"context"
	"time"

	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/notes/application"
	"github.com/zchelalo/neuraclinic-records/internal/shared/i18n"
	recordgrpc "github.com/zchelalo/neuraclinic-records/internal/shared/recordgrpc"
)

func (s *NoteService) Create(ctx context.Context, req *recordv1.NoteServiceCreateRequest) (*recordv1.NoteServiceCreateResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	patientID, err := recordgrpc.ParseID(req.GetPatientId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	appointmentID, err := recordgrpc.ParseOptionalID(req.AppointmentId)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	note, err := s.app.CreateNote(ctx, application.NoteCreateCommand{
		PsychologistID: psychologistID,
		PatientID:      patientID,
		AppointmentID:  appointmentID,
		Title:          req.Title,
		ContentHTML:    req.GetContentHtml(),
		ContentText:    req.GetContentText(),
	})
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	return &recordv1.NoteServiceCreateResponse{Note: recordgrpc.NoteToProto(note)}, nil
}

func (s *NoteService) List(ctx context.Context, req *recordv1.NoteServiceListRequest) (*recordv1.NoteServiceListResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	patientID, err := recordgrpc.ParseID(req.GetPatientId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	var startDate, endDate *time.Time
	if req.GetCreatedAtRange() != nil {
		if req.GetCreatedAtRange().GetStartDate() != nil {
			value := req.GetCreatedAtRange().GetStartDate().AsTime()
			startDate = &value
		}
		if req.GetCreatedAtRange().GetEndDate() != nil {
			value := req.GetCreatedAtRange().GetEndDate().AsTime()
			endDate = &value
		}
	}
	result, err := s.app.ListNotes(ctx, application.NoteListCommand{
		PsychologistID:            psychologistID,
		PatientID:                 patientID,
		Pagination:                recordgrpc.CursorPaginationFromProto(req.GetPagination()),
		StartDate:                 startDate,
		EndDate:                   endDate,
		WithAppointmentAssociated: req.GetWithAppointmentAssociated(),
		WithFilesAssociated:       req.GetWithFilesAssociated(),
		SearchQuery:               req.GetSearchQuery(),
	})
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	resp := make([]*recordv1.NoteSummary, 0, len(result.Notes))
	for _, note := range result.Notes {
		resp = append(resp, recordgrpc.NoteSummaryToProto(note))
	}
	return &recordv1.NoteServiceListResponse{
		Notes: resp,
		Meta:  recordgrpc.CursorMetaToProto(result.Meta),
	}, nil
}

func (s *NoteService) FindById(ctx context.Context, req *recordv1.NoteServiceFindByIdRequest) (*recordv1.NoteServiceFindByIdResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	note, err := s.app.FindNote(ctx, psychologistID, id)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	return &recordv1.NoteServiceFindByIdResponse{Note: recordgrpc.NoteToProto(note)}, nil
}

func (s *NoteService) Update(ctx context.Context, req *recordv1.NoteServiceUpdateRequest) (*recordv1.NoteServiceUpdateResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	appointmentID, err := recordgrpc.ParseOptionalID(req.AppointmentId)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	note, err := s.app.UpdateNote(ctx, application.NoteUpdateCommand{
		PsychologistID: psychologistID,
		ID:             id,
		AppointmentID:  appointmentID,
		Title:          req.Title,
		ContentHTML:    req.ContentHtml,
		ContentText:    req.ContentText,
	})
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	return &recordv1.NoteServiceUpdateResponse{Note: recordgrpc.NoteToProto(note)}, nil
}

func (s *NoteService) Delete(ctx context.Context, req *recordv1.NoteServiceDeleteRequest) (*recordv1.NoteServiceDeleteResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	if err := s.app.DeleteNote(ctx, psychologistID, id); err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	return &recordv1.NoteServiceDeleteResponse{Operation: recordgrpc.Operation(ctx, i18n.KeyNoteDeleted)}, nil
}
