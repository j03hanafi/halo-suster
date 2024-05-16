package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/id"
	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

type MedicalRepository struct {
	db *pgxpool.Pool
}

func NewMedicalRepository(db *pgxpool.Pool) *MedicalRepository {
	return &MedicalRepository{db: db}
}

func (r MedicalRepository) RecordPatient(ctx context.Context, patient *domain.Patient) error {
	callerInfo := "[MedicalRepository.RecordPatient]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	patient.CreatedAt = time.Now()

	insertQuery := `INSERT INTO patients (id, phone_number, name, birth_date, is_male, img_url, created_at) 
		VALUES (@id, @phone_number, @name, @birth_date, @is_male, @img_url, @created_at)`
	args := pgx.NamedArgs{
		"id":           patient.ID,
		"phone_number": patient.PhoneNumber,
		"name":         patient.Name,
		"birth_date":   patient.BirthDate,
		"is_male":      patient.Gender == domain.GenderMale,
		"img_url":      patient.ImgURL,
		"created_at":   patient.CreatedAt,
	}

	_, err := r.db.Exec(ctx, insertQuery, args)
	if err != nil {
		l.Error("failed to register user", zap.Error(err))

		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return new(domain.ErrDuplicatePatient)
		}

		return err
	}

	return nil
}

func (r MedicalRepository) GetPatients(
	ctx context.Context,
	filter *domain.FilterPatient,
	patients domain.Patients,
) (domain.Patients, error) {
	callerInfo := "[MedicalRepository.GetPatients]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	conditions, params := r.filterPatient(filter)
	getQuery := `SELECT id, phone_number, name, birth_date, is_male, created_at FROM patients` + conditions

	rows, err := r.db.Query(ctx, getQuery, params)
	if err != nil {
		l.Error("failed to get patients", zap.Error(err))
		return patients, err
	}

	dPatient := domain.PatientAcquire()
	defer domain.PatientRelease(dPatient)
	var isMale bool

	_, err = pgx.ForEachRow(
		rows,
		[]any{&dPatient.ID, &dPatient.PhoneNumber, &dPatient.Name, &dPatient.BirthDate, &isMale, &dPatient.CreatedAt},
		func() error {
			dPatient.Gender = domain.GenderMale
			if !isMale {
				dPatient.Gender = domain.GenderFemale
			}
			patients = append(patients, *dPatient)
			return nil
		},
	)
	if err != nil {
		l.Error("failed to get patients", zap.Error(err))
		return patients, err
	}

	return patients, nil
}

func (r MedicalRepository) filterPatient(filter *domain.FilterPatient) (string, pgx.NamedArgs) {
	const totalConditions = 3
	conditions, params := make([]string, 0, totalConditions), pgx.NamedArgs{}

	if filter.ID != "" {
		conditions = append(conditions, "id = @id")
		params["id"] = filter.ID
	}

	if filter.Name != "" {
		conditions = append(conditions, "name ILIKE @name")
		params["name"] = "%" + filter.Name + "%"
	}

	if filter.PhoneNumber != "" {
		conditions = append(conditions, "phone_number LIKE @phone_number")
		params["phone_number"] = filter.PhoneNumber + "%"
	}

	order := " ORDER BY created_at DESC"
	if filter.CreatedAt != "" && (filter.CreatedAt == "asc" || filter.CreatedAt == "desc") {
		order = " ORDER BY created_at " + filter.CreatedAt
	}

	const totalLimitOffset = 2
	limitOffset := make([]string, 0, totalLimitOffset)

	limitOffset = append(limitOffset, "LIMIT @limit")
	params["limit"] = 5
	if filter.Limit != 0 {
		params["limit"] = filter.Limit
	}

	if filter.Offset != 0 {
		limitOffset = append(limitOffset, "OFFSET @offset")
		params["offset"] = filter.Offset
	}

	queryConditions := ""
	if len(conditions) > 0 {
		queryConditions = " WHERE " + strings.Join(conditions, " AND ")
	}

	queryConditions += order

	if len(limitOffset) > 0 {
		queryConditions += " " + strings.Join(limitOffset, " ")
	}

	return queryConditions, params
}

func (r MedicalRepository) SaveMedicalRecord(ctx context.Context, record *domain.MedicalRecord) error {
	callerInfo := "[MedicalRepository.SaveMedicalRecord]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	record.ID = id.New()
	record.CreatedAt = time.Now()

	insertQuery := `INSERT INTO medical_records (
		id, patient_id, patient_phone_number, patient_name, patient_birth_date, patient_is_male, patient_img_url, 
                             symptoms, medications, staff_id, staff_nip, staff_name, created_at
			)
		SELECT 
    		@id, p.id, p.phone_number, p.name, p.birth_date, p.is_male, p.img_url, @symptoms, @medications, @staff_id, 
    		@staff_nip, @staff_name, @created_at
		FROM patients p WHERE p.id = @patient_id`
	args := pgx.NamedArgs{
		"id":          record.ID,
		"patient_id":  record.PatientID,
		"symptoms":    record.Symptoms,
		"medications": record.Medications,
		"staff_id":    record.StaffID,
		"staff_nip":   record.StaffNIP,
		"staff_name":  record.StaffName,
		"created_at":  record.CreatedAt,
	}

	result, err := r.db.Exec(ctx, insertQuery, args)
	if err != nil {
		l.Error("failed to save medical record", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return new(domain.ErrPatientNotFound)
	}

	return nil
}

func (r MedicalRepository) GetMedicalRecords(
	ctx context.Context,
	filter *domain.FilterMedicalRecord,
	records domain.MedicalRecords,
) (domain.MedicalRecords, error) {
	callerInfo := "[MedicalRepository.GetMedicalRecords]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	conditions, params := r.filterMedicalRecord(filter)
	getQuery := `SELECT id, patient_id, patient_phone_number, patient_name, patient_birth_date, patient_is_male, patient_img_url, symptoms, medications, staff_id, staff_nip, staff_name, created_at FROM medical_records
		` + conditions

	rows, err := r.db.Query(ctx, getQuery, params)
	if err != nil {
		l.Error("failed to get medical records", zap.Error(err))
		return records, err
	}

	dRecord := domain.MedicalRecordAcquire()
	defer domain.MedicalRecordRelease(dRecord)
	var isMale bool

	_, err = pgx.ForEachRow(
		rows,
		[]any{
			&dRecord.ID,
			&dRecord.PatientID,
			&dRecord.PatientPhoneNumber,
			&dRecord.PatientName,
			&dRecord.PatientBirthDate,
			&isMale,
			&dRecord.PatientImgURL,
			&dRecord.Symptoms,
			&dRecord.Medications,
			&dRecord.StaffID,
			&dRecord.StaffNIP,
			&dRecord.StaffName,
			&dRecord.CreatedAt,
		},
		func() error {
			dRecord.PatientGender = domain.GenderMale
			if !isMale {
				dRecord.PatientGender = domain.GenderFemale
			}
			records = append(records, *dRecord)
			return nil
		},
	)
	if err != nil {
		l.Error("failed to get medical records", zap.Error(err))
		return records, err
	}

	return records, nil
}

func (r MedicalRepository) filterMedicalRecord(filter *domain.FilterMedicalRecord) (string, pgx.NamedArgs) {
	const totalConditions = 3
	conditions, params := make([]string, 0, totalConditions), pgx.NamedArgs{}

	if filter.PatientID != "" {
		conditions = append(conditions, "patient_id = @patient_id")
		params["patient_id"] = filter.PatientID
	}

	if !id.IsZero(filter.StaffID) {
		conditions = append(conditions, "staff_id = @staff_id")
		params["staff_id"] = filter.StaffID
	}

	if filter.StaffNIP != "" {
		conditions = append(conditions, "staff_nip = @staff_nip")
		params["staff_nip"] = filter.StaffNIP
	}

	order := " ORDER BY created_at DESC"
	if filter.CreatedAt != "" && (filter.CreatedAt == "asc" || filter.CreatedAt == "desc") {
		order = " ORDER BY created_at " + filter.CreatedAt
	}

	const totalLimitOffset = 2
	limitOffset := make([]string, 0, totalLimitOffset)

	limitOffset = append(limitOffset, "LIMIT @limit")
	params["limit"] = 5
	if filter.Limit != 0 {
		params["limit"] = filter.Limit
	}

	if filter.Offset != 0 {
		limitOffset = append(limitOffset, "OFFSET @offset")
		params["offset"] = filter.Offset
	}

	queryConditions := ""
	if len(conditions) > 0 {
		queryConditions = " WHERE " + strings.Join(conditions, " AND ")
	}

	queryConditions += order

	if len(limitOffset) > 0 {
		queryConditions += " " + strings.Join(limitOffset, " ")
	}

	return queryConditions, params
}

var _ MedicalRepositoryContract = (*MedicalRepository)(nil)
