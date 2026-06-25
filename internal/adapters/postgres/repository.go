package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	recordv1 "github.com/zchelalo/neuraclinic-records/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-records/gen/go/shared/v1"
	recordsdb "github.com/zchelalo/neuraclinic-records/internal/db/sqlc/records"
	appointmentdomain "github.com/zchelalo/neuraclinic-records/internal/modules/appointments/domain"
	appointmentports "github.com/zchelalo/neuraclinic-records/internal/modules/appointments/ports"
	attachmentdomain "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/domain"
	attachmentports "github.com/zchelalo/neuraclinic-records/internal/modules/attachments/ports"
	familiogramdomain "github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/domain"
	familiogramports "github.com/zchelalo/neuraclinic-records/internal/modules/familiogram/ports"
	notedomain "github.com/zchelalo/neuraclinic-records/internal/modules/notes/domain"
	noteports "github.com/zchelalo/neuraclinic-records/internal/modules/notes/ports"
	patientdomain "github.com/zchelalo/neuraclinic-records/internal/modules/patients/domain"
	patientports "github.com/zchelalo/neuraclinic-records/internal/modules/patients/ports"
	pgutil "github.com/zchelalo/neuraclinic-records/internal/shared/postgresutil"
	appshared "github.com/zchelalo/neuraclinic-records/internal/shared/recordapp"
	recorderrors "github.com/zchelalo/neuraclinic-records/internal/shared/recorderrors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type Repository struct {
	db *pgxpool.Pool
	q  *recordsdb.Queries
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db, q: recordsdb.New(db)}
}

func (r *Repository) CreatePatient(ctx context.Context, patient patientdomain.PatientCreate) (patientdomain.PatientSummary, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return patientdomain.PatientSummary{}, err
	}
	defer rollback(ctx, tx)

	q := r.q.WithTx(tx)
	now := patient.Now.UTC()
	if err := q.CreateAddress(ctx, recordsdb.CreateAddressParams{
		ID:           pgutil.UUID(patient.AddressID),
		Country:      patient.Country,
		Province:     patient.Province,
		City:         patient.City,
		PostalCode:   patient.PostalCode,
		Neighborhood: patient.Neighborhood,
		Street:       patient.Street,
		StreetNumber: patient.StreetNumber,
		UnitNumber:   pgutil.OptionalText(patient.UnitNumber),
		CreatedAt:    pgutil.Timestamptz(now),
	}); err != nil {
		return patientdomain.PatientSummary{}, err
	}

	if err := q.CreatePatient(ctx, recordsdb.CreatePatientParams{
		ID:             pgutil.UUID(patient.ID),
		FirstName:      patient.FirstName,
		MiddleName:     pgutil.OptionalText(patient.MiddleName),
		FirstLastName:  patient.FirstLastName,
		SecondLastName: pgutil.OptionalText(patient.SecondLastName),
		BirthDate:      pgutil.Date(patient.BirthDate),
		BirthCountry:   patient.BirthCountry,
		BirthState:     patient.BirthProvince,
		BirthCity:      patient.BirthCity,
		Sex:            patient.Sex.String(),
		MaritalStatus:  patient.MaritalStatus.String(),
		Occupation:     pgutil.OptionalText(patient.Occupation),
		Religion:       pgutil.OptionalText(patient.Religion),
		Phone:          patient.Phone,
		Email:          patient.Email,
		AddressID:      pgutil.UUID(patient.AddressID),
		PsychologistID: pgutil.UUID(patient.PsychologistID),
		CreatedAt:      pgutil.Timestamptz(now),
	}); err != nil {
		return patientdomain.PatientSummary{}, err
	}

	if err := q.CreateFamiliogram(ctx, recordsdb.CreateFamiliogramParams{
		ID:        pgutil.UUID(patient.FamiliogramID),
		Data:      []byte(`{}`),
		PatientID: pgutil.UUID(patient.ID),
		CreatedAt: pgutil.Timestamptz(now),
	}); err != nil {
		return patientdomain.PatientSummary{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return patientdomain.PatientSummary{}, err
	}

	return patientdomain.PatientSummary{
		ID:             patient.ID,
		FirstName:      patient.FirstName,
		MiddleName:     patient.MiddleName,
		FirstLastName:  patient.FirstLastName,
		SecondLastName: patient.SecondLastName,
		BirthDate:      patient.BirthDate,
		Email:          patient.Email,
		Phone:          patient.Phone,
	}, nil
}

func (r *Repository) ListPatients(ctx context.Context, filter patientdomain.PatientListFilter) (patientdomain.PatientList, error) {
	queryLimit := filter.Pagination.QueryLimit()
	params := recordsdb.ListPatientsParams{
		PsychologistID:          pgutil.UUID(filter.PsychologistID),
		SearchQuery:             strings.TrimSpace(filter.SearchQuery),
		WithPendingAppointments: filter.WithPendingAppointments,
		WithNoAppointments:      filter.WithNoAppointments,
		EverHadAppointments:     filter.EverHadAppointments,
		AfterID:                 pgutil.OptionalUUID(filter.Pagination.AfterID),
		BeforeID:                pgutil.OptionalUUID(filter.Pagination.BeforeID),
		IsBackward:              filter.Pagination.IsBackward(),
		LimitCount:              queryLimit,
	}
	rows, err := r.q.ListPatients(ctx, params)
	if err != nil {
		return patientdomain.PatientList{}, err
	}
	if filter.Pagination.IsBackward() {
		rows = appshared.NormalizeBackwardCursorRows(rows, queryLimit)
	}

	patients := make([]patientdomain.PatientSummary, 0, len(rows))
	for _, row := range rows {
		patients = append(patients, patientSummaryFromRow(row))
	}
	page := appshared.BuildCursorPage(patients, filter.Pagination, func(patient patientdomain.PatientSummary) uuid.UUID {
		return patient.ID
	})
	return patientdomain.PatientList{
		Patients: page.Items,
		Meta:     page.Meta,
	}, nil
}

func (r *Repository) PatientByID(ctx context.Context, psychologistID, id uuid.UUID) (patientdomain.Patient, error) {
	row, err := r.q.GetPatientByID(ctx, recordsdb.GetPatientByIDParams{
		ID:             pgutil.UUID(id),
		PsychologistID: pgutil.UUID(psychologistID),
	})
	if err != nil {
		return patientdomain.Patient{}, mapNoRows(err)
	}
	return patientFromRow(row), nil
}

func (r *Repository) UpdatePatientIdentification(ctx context.Context, update patientdomain.PatientIdentificationUpdate) (patientdomain.Patient, error) {
	var sex *string
	if update.Sex != nil {
		value := update.Sex.String()
		sex = &value
	}
	var maritalStatus *string
	if update.MaritalStatus != nil {
		value := update.MaritalStatus.String()
		maritalStatus = &value
	}
	rows, err := r.q.UpdatePatientIdentification(ctx, recordsdb.UpdatePatientIdentificationParams{
		FirstName:      pgutil.OptionalText(update.FirstName),
		MiddleName:     pgutil.OptionalText(update.MiddleName),
		FirstLastName:  pgutil.OptionalText(update.FirstLastName),
		SecondLastName: pgutil.OptionalText(update.SecondLastName),
		BirthDate:      pgutil.OptionalDate(update.BirthDate),
		Sex:            pgutil.OptionalText(sex),
		BirthCountry:   pgutil.OptionalText(update.BirthCountry),
		BirthState:     pgutil.OptionalText(update.BirthProvince),
		BirthCity:      pgutil.OptionalText(update.BirthCity),
		Occupation:     pgutil.OptionalText(update.Occupation),
		MaritalStatus:  pgutil.OptionalText(maritalStatus),
		Religion:       pgutil.OptionalText(update.Religion),
		UpdatedAt:      pgutil.Timestamptz(update.Now.UTC()),
		ID:             pgutil.UUID(update.ID),
		PsychologistID: pgutil.UUID(update.PsychologistID),
	})
	if err != nil {
		return patientdomain.Patient{}, err
	}
	if rows == 0 {
		return patientdomain.Patient{}, recorderrors.ErrNotFound
	}
	return r.PatientByID(ctx, update.PsychologistID, update.ID)
}

func (r *Repository) UpdatePatientContact(ctx context.Context, update patientdomain.PatientContactUpdate) (patientdomain.Patient, error) {
	rows, err := r.q.UpdatePatientContact(ctx, recordsdb.UpdatePatientContactParams{
		Phone:          pgutil.OptionalText(update.Phone),
		Email:          pgutil.OptionalText(update.Email),
		UpdatedAt:      pgutil.Timestamptz(update.Now.UTC()),
		ID:             pgutil.UUID(update.ID),
		PsychologistID: pgutil.UUID(update.PsychologistID),
	})
	if err != nil {
		return patientdomain.Patient{}, err
	}
	if rows == 0 {
		return patientdomain.Patient{}, recorderrors.ErrNotFound
	}
	return r.PatientByID(ctx, update.PsychologistID, update.ID)
}

func (r *Repository) UpdatePatientAddress(ctx context.Context, update patientdomain.AddressUpdate) (patientdomain.Patient, error) {
	rows, err := r.q.UpdateAddressByPatientID(ctx, recordsdb.UpdateAddressByPatientIDParams{
		Country:        pgutil.OptionalText(update.Country),
		Province:       pgutil.OptionalText(update.Province),
		City:           pgutil.OptionalText(update.City),
		PostalCode:     pgutil.OptionalText(update.PostalCode),
		Neighborhood:   pgutil.OptionalText(update.Neighborhood),
		Street:         pgutil.OptionalText(update.Street),
		StreetNumber:   pgutil.OptionalText(update.StreetNumber),
		UnitNumber:     pgutil.OptionalText(update.UnitNumber),
		UpdatedAt:      pgutil.Timestamptz(update.Now.UTC()),
		PatientID:      pgutil.UUID(update.PatientID),
		PsychologistID: pgutil.UUID(update.PsychologistID),
	})
	if err != nil {
		return patientdomain.Patient{}, err
	}
	if rows == 0 {
		return patientdomain.Patient{}, recorderrors.ErrNotFound
	}
	return r.PatientByID(ctx, update.PsychologistID, update.PatientID)
}

func (r *Repository) PatientExists(ctx context.Context, psychologistID, id uuid.UUID) (bool, error) {
	return r.q.PatientExists(ctx, recordsdb.PatientExistsParams{
		ID:             pgutil.UUID(id),
		PsychologistID: pgutil.UUID(psychologistID),
	})
}

func (r *Repository) FamiliogramByPatientID(ctx context.Context, psychologistID, patientID uuid.UUID) (familiogramdomain.Familiogram, error) {
	row, err := r.q.GetFamiliogramByPatientID(ctx, recordsdb.GetFamiliogramByPatientIDParams{
		PatientID:      pgutil.UUID(patientID),
		PsychologistID: pgutil.UUID(psychologistID),
	})
	if err != nil {
		return familiogramdomain.Familiogram{}, mapNoRows(err)
	}
	return familiogramFromRow(row)
}

func (r *Repository) UpdateFamiliogram(ctx context.Context, psychologistID, id uuid.UUID, data *structpb.Struct, now time.Time) (familiogramdomain.Familiogram, error) {
	raw, err := marshalStruct(data)
	if err != nil {
		return familiogramdomain.Familiogram{}, err
	}
	row, err := r.q.UpdateFamiliogram(ctx, recordsdb.UpdateFamiliogramParams{
		ID:             pgutil.UUID(id),
		PsychologistID: pgutil.UUID(psychologistID),
		Data:           raw,
		UpdatedAt:      pgutil.Timestamptz(now.UTC()),
	})
	if err != nil {
		return familiogramdomain.Familiogram{}, mapNoRows(err)
	}
	return familiogramFromRow(row)
}

func (r *Repository) CreateAppointment(ctx context.Context, appointment appointmentdomain.AppointmentCreate) (appointmentdomain.Appointment, error) {
	row, err := r.q.CreateAppointment(ctx, createAppointmentParams(appointment))
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	return appointmentFromRow(row), nil
}

func (r *Repository) AppointmentByID(ctx context.Context, psychologistID, id uuid.UUID) (appointmentdomain.Appointment, error) {
	row, err := r.q.GetAppointmentByID(ctx, recordsdb.GetAppointmentByIDParams{
		ID:             pgutil.UUID(id),
		PsychologistID: pgutil.UUID(psychologistID),
	})
	if err != nil {
		return appointmentdomain.Appointment{}, mapNoRows(err)
	}
	return appointmentFromRow(row), nil
}

func (r *Repository) ListAppointments(ctx context.Context, filter appointmentdomain.AppointmentListFilter) (appointmentdomain.AppointmentList, error) {
	queryLimit := filter.Pagination.QueryLimit()
	statuses := make([]string, 0, len(filter.Statuses))
	for _, status := range filter.Statuses {
		if status != sharedv1.AppointmentStatus_APPOINTMENT_STATUS_UNSPECIFIED {
			statuses = append(statuses, status.String())
		}
	}
	params := recordsdb.ListAppointmentsParams{
		PsychologistID: pgutil.UUID(filter.PsychologistID),
		PatientID:      pgutil.OptionalUUID(filter.PatientID),
		StartDate:      pgutil.OptionalTimestamptz(filter.StartDate),
		EndDate:        pgutil.OptionalTimestamptz(filter.EndDate),
		Statuses:       statuses,
		AfterID:        pgutil.OptionalUUID(filter.Pagination.AfterID),
		BeforeID:       pgutil.OptionalUUID(filter.Pagination.BeforeID),
		IsBackward:     filter.Pagination.IsBackward(),
		LimitCount:     queryLimit,
	}
	rows, err := r.q.ListAppointments(ctx, params)
	if err != nil {
		return appointmentdomain.AppointmentList{}, err
	}
	if filter.Pagination.IsBackward() {
		rows = appshared.NormalizeBackwardCursorRows(rows, queryLimit)
	}

	appointments := make([]appointmentdomain.Appointment, 0, len(rows))
	for _, row := range rows {
		appointments = append(appointments, appointmentFromRow(row))
	}
	page := appshared.BuildCursorPage(appointments, filter.Pagination, func(appointment appointmentdomain.Appointment) uuid.UUID {
		return appointment.ID
	})
	return appointmentdomain.AppointmentList{
		Appointments: page.Items,
		Meta:         page.Meta,
	}, nil
}

func (r *Repository) RescheduleAppointment(ctx context.Context, psychologistID, originalID uuid.UUID, appointment appointmentdomain.AppointmentCreate, now time.Time) (appointmentdomain.Appointment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	defer rollback(ctx, tx)

	q := r.q.WithTx(tx)
	_, err = q.UpdateAppointmentStatus(ctx, recordsdb.UpdateAppointmentStatusParams{
		ID:                pgutil.UUID(originalID),
		PsychologistID:    pgutil.UUID(psychologistID),
		Status:            sharedv1.AppointmentStatus_APPOINTMENT_STATUS_RESCHEDULED.String(),
		CancelledByUserID: pgtype.UUID{Valid: false},
		UpdatedAt:         pgutil.Timestamptz(now.UTC()),
	})
	if err != nil {
		return appointmentdomain.Appointment{}, mapNoRows(err)
	}

	row, err := q.CreateAppointment(ctx, createAppointmentParams(appointment))
	if err != nil {
		return appointmentdomain.Appointment{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return appointmentdomain.Appointment{}, err
	}
	return appointmentFromRow(row), nil
}

func (r *Repository) UpdateAppointmentStatus(ctx context.Context, psychologistID, id uuid.UUID, status sharedv1.AppointmentStatus, cancelledByUserID *uuid.UUID, now time.Time) (appointmentdomain.Appointment, error) {
	row, err := r.q.UpdateAppointmentStatus(ctx, recordsdb.UpdateAppointmentStatusParams{
		ID:                pgutil.UUID(id),
		PsychologistID:    pgutil.UUID(psychologistID),
		Status:            status.String(),
		CancelledByUserID: pgutil.OptionalUUID(cancelledByUserID),
		UpdatedAt:         pgutil.Timestamptz(now.UTC()),
	})
	if err != nil {
		return appointmentdomain.Appointment{}, mapNoRows(err)
	}
	return appointmentFromRow(row), nil
}

func (r *Repository) AppointmentBelongsToPatient(ctx context.Context, psychologistID, appointmentID, patientID uuid.UUID) (bool, error) {
	return r.q.AppointmentBelongsToPatient(ctx, recordsdb.AppointmentBelongsToPatientParams{
		ID:             pgutil.UUID(appointmentID),
		PatientID:      pgutil.UUID(patientID),
		PsychologistID: pgutil.UUID(psychologistID),
	})
}

func (r *Repository) CreateNote(ctx context.Context, note notedomain.NoteCreate) (notedomain.Note, error) {
	row, err := r.q.CreateNote(ctx, recordsdb.CreateNoteParams{
		ID:            pgutil.UUID(note.ID),
		PatientID:     pgutil.UUID(note.PatientID),
		AppointmentID: pgutil.OptionalUUID(note.AppointmentID),
		Title:         pgutil.OptionalText(note.Title),
		ContentHtml:   note.ContentHTML,
		ContentText:   note.ContentText,
		CreatedAt:     pgutil.Timestamptz(note.Now.UTC()),
	})
	if err != nil {
		return notedomain.Note{}, err
	}
	return noteFromRow(row), nil
}

func (r *Repository) NoteByID(ctx context.Context, psychologistID, id uuid.UUID) (notedomain.Note, error) {
	row, err := r.q.GetNoteByID(ctx, recordsdb.GetNoteByIDParams{
		ID:             pgutil.UUID(id),
		PsychologistID: pgutil.UUID(psychologistID),
	})
	if err != nil {
		return notedomain.Note{}, mapNoRows(err)
	}
	return noteFromRow(row), nil
}

func (r *Repository) ListNotes(ctx context.Context, filter notedomain.NoteListFilter) (notedomain.NoteList, error) {
	queryLimit := filter.Pagination.QueryLimit()
	rows, err := r.q.ListNotes(ctx, recordsdb.ListNotesParams{
		PsychologistID:            pgutil.UUID(filter.PsychologistID),
		PatientID:                 pgutil.UUID(filter.PatientID),
		StartDate:                 pgutil.OptionalTimestamptz(filter.StartDate),
		EndDate:                   pgutil.OptionalTimestamptz(filter.EndDate),
		WithAppointmentAssociated: filter.WithAppointmentAssociated,
		WithFilesAssociated:       filter.WithFilesAssociated,
		SearchQuery:               strings.TrimSpace(filter.SearchQuery),
		AfterID:                   pgutil.OptionalUUID(filter.Pagination.AfterID),
		BeforeID:                  pgutil.OptionalUUID(filter.Pagination.BeforeID),
		IsBackward:                filter.Pagination.IsBackward(),
		LimitCount:                queryLimit,
	})
	if err != nil {
		return notedomain.NoteList{}, err
	}
	if filter.Pagination.IsBackward() {
		rows = appshared.NormalizeBackwardCursorRows(rows, queryLimit)
	}
	notes := make([]notedomain.NoteSummary, 0, len(rows))
	for _, row := range rows {
		notes = append(notes, noteSummaryFromRow(row))
	}
	page := appshared.BuildCursorPage(notes, filter.Pagination, func(note notedomain.NoteSummary) uuid.UUID {
		return note.ID
	})
	return notedomain.NoteList{
		Notes: page.Items,
		Meta:  page.Meta,
	}, nil
}

func (r *Repository) UpdateNote(ctx context.Context, update notedomain.NoteUpdate) (notedomain.Note, error) {
	rows, err := r.q.UpdateNote(ctx, recordsdb.UpdateNoteParams{
		AppointmentID:  pgutil.OptionalUUID(update.AppointmentID),
		Title:          pgutil.OptionalText(update.Title),
		ContentHtml:    pgutil.OptionalText(update.ContentHTML),
		ContentText:    pgutil.OptionalText(update.ContentText),
		UpdatedAt:      pgutil.Timestamptz(update.Now.UTC()),
		ID:             pgutil.UUID(update.ID),
		PsychologistID: pgutil.UUID(update.PsychologistID),
	})
	if err != nil {
		return notedomain.Note{}, err
	}
	if rows == 0 {
		return notedomain.Note{}, recorderrors.ErrNotFound
	}
	return r.NoteByID(ctx, update.PsychologistID, update.ID)
}

func (r *Repository) DeleteNote(ctx context.Context, psychologistID, id uuid.UUID, now time.Time) error {
	rows, err := r.q.DeleteNote(ctx, recordsdb.DeleteNoteParams{
		ID:             pgutil.UUID(id),
		PsychologistID: pgutil.UUID(psychologistID),
		DeletedAt:      pgutil.Timestamptz(now.UTC()),
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return recorderrors.ErrNotFound
	}
	return nil
}

func (r *Repository) NoteBelongsToPatient(ctx context.Context, psychologistID, noteID, patientID uuid.UUID) (bool, error) {
	return r.q.NoteBelongsToPatient(ctx, recordsdb.NoteBelongsToPatientParams{
		ID:             pgutil.UUID(noteID),
		PatientID:      pgutil.UUID(patientID),
		PsychologistID: pgutil.UUID(psychologistID),
	})
}

func (r *Repository) CreateAttachment(ctx context.Context, attachment attachmentdomain.AttachmentCreate, fileID uuid.UUID) (attachmentdomain.Attachment, error) {
	row, err := r.q.CreateAttachment(ctx, recordsdb.CreateAttachmentParams{
		ID:           pgutil.UUID(attachment.ID),
		FileID:       pgutil.UUID(fileID),
		MimeType:     attachment.MimeType,
		PatientID:    pgutil.UUID(attachment.PatientID),
		NoteID:       pgutil.OptionalUUID(attachment.NoteID),
		UploadStatus: sharedv1.FileStatus_FILE_STATUS_UPLOADING.String(),
		CreatedAt:    pgutil.Timestamptz(attachment.Now.UTC()),
	})
	if err != nil {
		return attachmentdomain.Attachment{}, err
	}
	return attachmentFromRow(row), nil
}

func (r *Repository) AttachmentByID(ctx context.Context, psychologistID, id uuid.UUID) (attachmentdomain.Attachment, error) {
	row, err := r.q.GetAttachmentByID(ctx, recordsdb.GetAttachmentByIDParams{
		ID:             pgutil.UUID(id),
		PsychologistID: pgutil.UUID(psychologistID),
	})
	if err != nil {
		return attachmentdomain.Attachment{}, mapNoRows(err)
	}
	return attachmentFromRow(row), nil
}

func (r *Repository) ListAttachments(ctx context.Context, filter attachmentdomain.AttachmentListFilter) (attachmentdomain.AttachmentList, error) {
	queryLimit := filter.Pagination.QueryLimit()
	rows, err := r.q.ListAttachments(ctx, recordsdb.ListAttachmentsParams{
		PsychologistID: pgutil.UUID(filter.PsychologistID),
		PatientID:      pgutil.UUID(filter.PatientID),
		NoteID:         pgutil.OptionalUUID(filter.NoteID),
		AfterID:        pgutil.OptionalUUID(filter.Pagination.AfterID),
		BeforeID:       pgutil.OptionalUUID(filter.Pagination.BeforeID),
		IsBackward:     filter.Pagination.IsBackward(),
		LimitCount:     queryLimit,
	})
	if err != nil {
		return attachmentdomain.AttachmentList{}, err
	}
	if filter.Pagination.IsBackward() {
		rows = appshared.NormalizeBackwardCursorRows(rows, queryLimit)
	}
	attachments := make([]attachmentdomain.Attachment, 0, len(rows))
	for _, row := range rows {
		attachments = append(attachments, attachmentFromRow(row))
	}
	page := appshared.BuildCursorPage(attachments, filter.Pagination, func(attachment attachmentdomain.Attachment) uuid.UUID {
		return attachment.ID
	})
	return attachmentdomain.AttachmentList{
		Attachments: page.Items,
		Meta:        page.Meta,
	}, nil
}

func (r *Repository) UpdateAttachmentUploadStatusByFileID(ctx context.Context, fileID uuid.UUID, status sharedv1.FileStatus, now time.Time) (attachmentdomain.Attachment, error) {
	row, err := r.q.UpdateAttachmentUploadStatusByFileID(ctx, recordsdb.UpdateAttachmentUploadStatusByFileIDParams{
		FileID:       pgutil.UUID(fileID),
		UploadStatus: status.String(),
		UpdatedAt:    pgutil.Timestamptz(now.UTC()),
	})
	if err != nil {
		return attachmentdomain.Attachment{}, mapNoRows(err)
	}
	return attachmentFromRow(row), nil
}

func (r *Repository) DeleteAttachment(ctx context.Context, psychologistID, id uuid.UUID, now time.Time) error {
	rows, err := r.q.DeleteAttachment(ctx, recordsdb.DeleteAttachmentParams{
		ID:             pgutil.UUID(id),
		PsychologistID: pgutil.UUID(psychologistID),
		DeletedAt:      pgutil.Timestamptz(now.UTC()),
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return recorderrors.ErrNotFound
	}
	return nil
}

func createAppointmentParams(appointment appointmentdomain.AppointmentCreate) recordsdb.CreateAppointmentParams {
	return recordsdb.CreateAppointmentParams{
		ID:                           pgutil.UUID(appointment.ID),
		StartTime:                    pgutil.Timestamptz(appointment.StartTime.UTC()),
		EndTime:                      pgutil.Timestamptz(appointment.EndTime.UTC()),
		Reason:                       appointment.Reason,
		Status:                       appointment.Status.String(),
		PatientID:                    pgutil.UUID(appointment.PatientID),
		CancelledByUserID:            pgutil.OptionalUUID(appointment.CancelledByUserID),
		RescheduledFromAppointmentID: pgutil.OptionalUUID(appointment.RescheduledFromAppointmentID),
		CreatedAt:                    pgutil.Timestamptz(appointment.Now.UTC()),
	}
}

func patientFromRow(row recordsdb.GetPatientByIDRow) patientdomain.Patient {
	return patientdomain.Patient{
		ID:             pgutil.UUIDValue(row.ID),
		FirstName:      row.FirstName,
		MiddleName:     pgutil.TextPtr(row.MiddleName),
		FirstLastName:  row.FirstLastName,
		SecondLastName: pgutil.TextPtr(row.SecondLastName),
		BirthDate:      row.BirthDate.Time,
		BirthCountry:   row.BirthCountry,
		BirthProvince:  row.BirthState,
		BirthCity:      row.BirthCity,
		Sex:            parseSex(row.Sex),
		MaritalStatus:  parseMaritalStatus(row.MaritalStatus),
		Occupation:     pgutil.TextPtr(row.Occupation),
		Religion:       pgutil.TextPtr(row.Religion),
		Phone:          row.Phone,
		Email:          row.Email,
		PsychologistID: pgutil.UUIDValue(row.PsychologistID),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
		DeletedAt:      pgutil.TimestamptzPtr(row.DeletedAt),
		Address: patientdomain.Address{
			ID:           pgutil.UUIDValue(row.AddressID),
			Country:      row.AddressCountry,
			Province:     row.AddressProvince,
			City:         row.AddressCity,
			PostalCode:   row.AddressPostalCode,
			Neighborhood: row.AddressNeighborhood,
			Street:       row.AddressStreet,
			StreetNumber: row.AddressStreetNumber,
			UnitNumber:   pgutil.TextPtr(row.AddressUnitNumber),
			CreatedAt:    row.AddressCreatedAt.Time,
			UpdatedAt:    row.AddressUpdatedAt.Time,
			DeletedAt:    pgutil.TimestamptzPtr(row.AddressDeletedAt),
		},
	}
}

func patientSummaryFromRow(row recordsdb.ListPatientsRow) patientdomain.PatientSummary {
	return patientdomain.PatientSummary{
		ID:             pgutil.UUIDValue(row.ID),
		FirstName:      row.FirstName,
		MiddleName:     pgutil.TextPtr(row.MiddleName),
		FirstLastName:  row.FirstLastName,
		SecondLastName: pgutil.TextPtr(row.SecondLastName),
		BirthDate:      row.BirthDate.Time,
		Email:          row.Email,
		Phone:          row.Phone,
	}
}

func familiogramFromRow(row recordsdb.Familiogram) (familiogramdomain.Familiogram, error) {
	data := &structpb.Struct{}
	if len(row.Data) > 0 {
		if err := protojson.Unmarshal(row.Data, data); err != nil {
			return familiogramdomain.Familiogram{}, fmt.Errorf("unmarshal familiogram data: %w", err)
		}
	}
	return familiogramdomain.Familiogram{
		ID:        pgutil.UUIDValue(row.ID),
		Data:      data,
		PatientID: pgutil.UUIDValue(row.PatientID),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func appointmentFromRow(row recordsdb.Appointment) appointmentdomain.Appointment {
	return appointmentdomain.Appointment{
		ID:                           pgutil.UUIDValue(row.ID),
		StartTime:                    row.StartTime.Time,
		EndTime:                      row.EndTime.Time,
		Reason:                       row.Reason,
		Status:                       parseAppointmentStatus(row.Status),
		PatientID:                    pgutil.UUIDValue(row.PatientID),
		CancelledByUserID:            pgutil.UUIDPtr(row.CancelledByUserID),
		RescheduledFromAppointmentID: pgutil.UUIDPtr(row.RescheduledFromAppointmentID),
		CreatedAt:                    row.CreatedAt.Time,
		UpdatedAt:                    row.UpdatedAt.Time,
	}
}

func noteFromRow(row recordsdb.Note) notedomain.Note {
	return notedomain.Note{
		ID:            pgutil.UUIDValue(row.ID),
		PatientID:     pgutil.UUIDValue(row.PatientID),
		AppointmentID: pgutil.UUIDPtr(row.AppointmentID),
		Title:         pgutil.TextPtr(row.Title),
		ContentHTML:   row.ContentHtml,
		ContentText:   row.ContentText,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		DeletedAt:     pgutil.TimestamptzPtr(row.DeletedAt),
	}
}

func noteSummaryFromRow(row recordsdb.ListNotesRow) notedomain.NoteSummary {
	return notedomain.NoteSummary{
		ID:            pgutil.UUIDValue(row.ID),
		PatientID:     pgutil.UUIDValue(row.PatientID),
		AppointmentID: pgutil.UUIDPtr(row.AppointmentID),
		Title:         pgutil.TextPtr(row.Title),
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
		DeletedAt:     pgutil.TimestamptzPtr(row.DeletedAt),
	}
}

func attachmentFromRow(row recordsdb.Attachment) attachmentdomain.Attachment {
	return attachmentdomain.Attachment{
		ID:           pgutil.UUIDValue(row.ID),
		FileID:       pgutil.UUIDValue(row.FileID),
		MimeType:     row.MimeType,
		UploadStatus: parseFileStatus(row.UploadStatus),
		PatientID:    pgutil.UUIDValue(row.PatientID),
		NoteID:       pgutil.UUIDPtr(row.NoteID),
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		DeletedAt:    pgutil.TimestamptzPtr(row.DeletedAt),
	}
}

func parseFileStatus(value string) sharedv1.FileStatus {
	if parsed, ok := sharedv1.FileStatus_value[value]; ok {
		return sharedv1.FileStatus(parsed)
	}
	return sharedv1.FileStatus_FILE_STATUS_UNSPECIFIED
}

func parseSex(value string) recordv1.Sex {
	if parsed, ok := recordv1.Sex_value[value]; ok {
		return recordv1.Sex(parsed)
	}
	return recordv1.Sex_SEX_UNSPECIFIED
}

func parseMaritalStatus(value string) recordv1.MaritalStatus {
	if parsed, ok := recordv1.MaritalStatus_value[value]; ok {
		return recordv1.MaritalStatus(parsed)
	}
	return recordv1.MaritalStatus_MARITAL_STATUS_UNSPECIFIED
}

func parseAppointmentStatus(value string) sharedv1.AppointmentStatus {
	if parsed, ok := sharedv1.AppointmentStatus_value[value]; ok {
		return sharedv1.AppointmentStatus(parsed)
	}
	return sharedv1.AppointmentStatus_APPOINTMENT_STATUS_UNSPECIFIED
}

func marshalStruct(data *structpb.Struct) ([]byte, error) {
	if data == nil {
		return []byte(`{}`), nil
	}
	raw, err := protojson.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal familiogram data: %w", err)
	}
	if !json.Valid(raw) {
		return nil, fmt.Errorf("familiogram data is not valid json")
	}
	return raw, nil
}

func mapNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return recorderrors.ErrNotFound
	}
	return err
}

func rollback(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}

var _ patientports.Repository = (*Repository)(nil)
var _ appointmentports.Repository = (*Repository)(nil)
var _ noteports.Repository = (*Repository)(nil)
var _ familiogramports.Repository = (*Repository)(nil)
var _ attachmentports.Repository = (*Repository)(nil)
