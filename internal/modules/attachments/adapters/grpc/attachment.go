package grpc

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-records/internal/modules/attachments/application"
	"github.com/zchelalo/neuraclinic-records/internal/shared/i18n"
	recordgrpc "github.com/zchelalo/neuraclinic-records/internal/shared/recordgrpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AttachmentService) Create(ctx context.Context, req *recordv1.AttachmentServiceCreateRequest) (*recordv1.AttachmentServiceCreateResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	patientID, err := recordgrpc.ParseID(req.GetPatientId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	noteID, err := recordgrpc.ParseOptionalID(req.NoteId)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	result, err := s.app.CreateAttachment(ctx, application.AttachmentCreateCommand{
		PsychologistID: psychologistID,
		UserID:         recordgrpc.UserID(ctx),
		PatientID:      patientID,
		NoteID:         noteID,
		OriginalName:   req.GetOriginalName(),
		MimeType:       req.GetMimeType(),
		SizeBytes:      req.GetSizeBytes(),
	})
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	return &recordv1.AttachmentServiceCreateResponse{
		Id:        result.ID.String(),
		FileId:    result.FileID.String(),
		UploadUrl: result.UploadURL,
		ExpiresAt: timestamppb.New(result.ExpiresAt),
	}, nil
}

func (s *AttachmentService) List(ctx context.Context, req *recordv1.AttachmentServiceListRequest) (*recordv1.AttachmentServiceListResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	patientID, err := recordgrpc.ParseID(req.GetPatientId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	noteID, err := recordgrpc.ParseOptionalID(req.NoteId)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	result, err := s.app.ListAttachments(ctx, application.AttachmentListCommand{
		PsychologistID: psychologistID,
		PatientID:      patientID,
		NoteID:         noteID,
		Pagination:     recordgrpc.CursorPaginationFromProto(req.GetPagination()),
	})
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	resp := make([]*recordv1.Attachment, 0, len(result.Attachments))
	for _, attachment := range result.Attachments {
		resp = append(resp, recordgrpc.AttachmentToProto(attachment))
	}
	return &recordv1.AttachmentServiceListResponse{
		Attachments: resp,
		Meta:        recordgrpc.CursorMetaToProto(result.Meta),
	}, nil
}

func (s *AttachmentService) FindById(ctx context.Context, req *recordv1.AttachmentServiceFindByIdRequest) (*recordv1.AttachmentServiceFindByIdResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	attachment, err := s.app.FindAttachment(ctx, psychologistID, id)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	return &recordv1.AttachmentServiceFindByIdResponse{Attachment: recordgrpc.AttachmentToProto(attachment)}, nil
}

func (s *AttachmentService) Delete(ctx context.Context, req *recordv1.AttachmentServiceDeleteRequest) (*recordv1.AttachmentServiceDeleteResponse, error) {
	psychologistID, err := recordgrpc.PsychologistID(ctx)
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	id, err := recordgrpc.ParseID(req.GetId())
	if err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	if err := s.app.DeleteAttachment(ctx, psychologistID, id); err != nil {
		return nil, recordgrpc.MapError(ctx, err)
	}
	return &recordv1.AttachmentServiceDeleteResponse{Operation: recordgrpc.Operation(ctx, i18n.KeyAttachmentDeleted)}, nil
}
