package recordgrpc

import (
	"fmt"
	"time"

	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	appointmentdomain "github.com/zchelalo/neuraclinic-records/internal/modules/appointments/domain"
	attachmentdomain "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	familiogramdomain "github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/domain"
	notedomain "github.com/zchelalo/neuraclinic-records/internal/modules/notes/domain"
	patientdomain "github.com/zchelalo/neuraclinic-records/internal/modules/patients/domain"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	date "google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimeFromProtoDate(value *date.Date) (time.Time, error) {
	if value == nil || value.GetYear() <= 0 || value.GetMonth() <= 0 || value.GetDay() <= 0 {
		return time.Time{}, fmt.Errorf("invalid date")
	}
	return time.Date(int(value.GetYear()), time.Month(value.GetMonth()), int(value.GetDay()), 0, 0, 0, 0, time.UTC), nil
}

func protoDateFromTime(value time.Time) *date.Date {
	value = value.UTC()
	return &date.Date{
		Year:  int32(value.Year()),
		Month: int32(value.Month()),
		Day:   int32(value.Day()),
	}
}

func PatientSummaryToProto(patient patientdomain.PatientSummary) *recordv1.PatientSummary {
	return &recordv1.PatientSummary{
		Id:             patient.ID.String(),
		FirstName:      patient.FirstName,
		MiddleName:     patient.MiddleName,
		FirstLastName:  patient.FirstLastName,
		SecondLastName: patient.SecondLastName,
		BirthDate:      patient.BirthDate.Format(time.DateOnly),
		Email:          patient.Email,
		Phone:          patient.Phone,
	}
}

func PatientToProto(patient patientdomain.Patient) *recordv1.Patient {
	resp := &recordv1.Patient{
		Id:             patient.ID.String(),
		FirstName:      patient.FirstName,
		MiddleName:     patient.MiddleName,
		FirstLastName:  patient.FirstLastName,
		SecondLastName: patient.SecondLastName,
		BirthDate:      protoDateFromTime(patient.BirthDate),
		BirthCountry:   patient.BirthCountry,
		BirthProvince:  patient.BirthProvince,
		BirthCity:      patient.BirthCity,
		Sex:            patient.Sex,
		MaritalStatus:  patient.MaritalStatus,
		Occupation:     patient.Occupation,
		Religion:       patient.Religion,
		Phone:          patient.Phone,
		Email:          patient.Email,
		Address:        addressToProto(patient.Address),
		PsychologistId: patient.PsychologistID.String(),
		CreatedAt:      timestamppb.New(patient.CreatedAt),
		UpdatedAt:      timestamppb.New(patient.UpdatedAt),
	}
	if patient.DeletedAt != nil {
		resp.DeletedAt = timestamppb.New(*patient.DeletedAt)
	}
	return resp
}

func addressToProto(address patientdomain.Address) *recordv1.Address {
	resp := &recordv1.Address{
		Id:           address.ID.String(),
		Country:      address.Country,
		Province:     address.Province,
		City:         address.City,
		PostalCode:   address.PostalCode,
		Neighborhood: address.Neighborhood,
		Street:       address.Street,
		StreetNumber: address.StreetNumber,
		UnitNumber:   address.UnitNumber,
		CreatedAt:    timestamppb.New(address.CreatedAt),
		UpdatedAt:    timestamppb.New(address.UpdatedAt),
	}
	if address.DeletedAt != nil {
		resp.DeletedAt = timestamppb.New(*address.DeletedAt)
	}
	return resp
}

func AppointmentToProto(appointment appointmentdomain.Appointment) *recordv1.Appointment {
	resp := &recordv1.Appointment{
		Id:        appointment.ID.String(),
		StartTime: timestamppb.New(appointment.StartTime),
		EndTime:   timestamppb.New(appointment.EndTime),
		Reason:    appointment.Reason,
		Status:    appointment.Status,
		PatientId: appointment.PatientID.String(),
		CreatedAt: timestamppb.New(appointment.CreatedAt),
		UpdatedAt: timestamppb.New(appointment.UpdatedAt),
	}
	if appointment.CancelledByUserID != nil {
		value := appointment.CancelledByUserID.String()
		resp.CancelledByUserId = &value
	}
	if appointment.RescheduledFromAppointmentID != nil {
		value := appointment.RescheduledFromAppointmentID.String()
		resp.RescheduledFromAppointmentId = &value
	}
	return resp
}

func NoteToProto(note notedomain.Note) *recordv1.Note {
	resp := &recordv1.Note{
		Id:          note.ID.String(),
		PatientId:   note.PatientID.String(),
		Title:       note.Title,
		ContentHtml: note.ContentHTML,
		ContentText: note.ContentText,
		CreatedAt:   timestamppb.New(note.CreatedAt),
		UpdatedAt:   timestamppb.New(note.UpdatedAt),
	}
	if note.AppointmentID != nil {
		value := note.AppointmentID.String()
		resp.AppointmentId = &value
	}
	if note.DeletedAt != nil {
		resp.DeletedAt = timestamppb.New(*note.DeletedAt)
	}
	return resp
}

func NoteSummaryToProto(note notedomain.NoteSummary) *recordv1.NoteSummary {
	resp := &recordv1.NoteSummary{
		Id:        note.ID.String(),
		PatientId: note.PatientID.String(),
		Title:     note.Title,
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: timestamppb.New(note.UpdatedAt),
	}
	if note.AppointmentID != nil {
		value := note.AppointmentID.String()
		resp.AppointmentId = &value
	}
	if note.DeletedAt != nil {
		resp.DeletedAt = timestamppb.New(*note.DeletedAt)
	}
	return resp
}

func FamiliogramToProto(familiogram familiogramdomain.Familiogram) *recordv1.Familiogram {
	return &recordv1.Familiogram{
		Id:        familiogram.ID.String(),
		Data:      familiogram.Data,
		PatientId: familiogram.PatientID.String(),
		CreatedAt: timestamppb.New(familiogram.CreatedAt),
		UpdatedAt: timestamppb.New(familiogram.UpdatedAt),
	}
}

func AttachmentToProto(attachment attachmentdomain.Attachment) *recordv1.Attachment {
	resp := &recordv1.Attachment{
		Id:           attachment.ID.String(),
		FileId:       attachment.FileID.String(),
		MimeType:     attachment.MimeType,
		UploadStatus: attachment.UploadStatus,
		PatientId:    attachment.PatientID.String(),
		CreatedAt:    timestamppb.New(attachment.CreatedAt),
		UpdatedAt:    timestamppb.New(attachment.UpdatedAt),
	}
	if attachment.DownloadURL != nil {
		resp.DownloadUrl = attachment.DownloadURL
	}
	if attachment.ExpiresAt != nil {
		resp.ExpiresAt = timestamppb.New(*attachment.ExpiresAt)
	}
	if attachment.NoteID != nil {
		value := attachment.NoteID.String()
		resp.NoteId = &value
	}
	if attachment.DeletedAt != nil {
		resp.DeletedAt = timestamppb.New(*attachment.DeletedAt)
	}
	return resp
}

func CursorMetaToProto(meta appshared.CursorMeta) *sharedv1.CursorMeta {
	return &sharedv1.CursorMeta{
		NextCursor: meta.NextCursor,
		PrevCursor: meta.PrevCursor,
		Limit:      meta.Limit,
	}
}

func CursorPaginationFromProto(p *sharedv1.CursorPagination) appshared.CursorPagination {
	if p == nil {
		return appshared.CursorPagination{}
	}
	return appshared.CursorPagination{
		AfterCursor:  p.AfterCursor,
		BeforeCursor: p.BeforeCursor,
		Limit:        p.GetLimit(),
	}
}

func Operation(message string) *sharedv1.OperationResponse {
	return &sharedv1.OperationResponse{Message: message}
}
