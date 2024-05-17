package service

import (
	"context"

	"github.com/j03hanafi/halo-suster/internal/domain"
)

type MedicalServiceContract interface {
	RecordPatient(ctx context.Context, patient *domain.Patient) error
	GetPatients(ctx context.Context, filter *domain.FilterPatient, patients domain.Patients) (domain.Patients, error)
	SaveMedicalRecord(ctx context.Context, record *domain.MedicalRecord, user *domain.User) error
	GetMedicalRecords(
		ctx context.Context,
		filter *domain.FilterMedicalRecord,
		records domain.MedicalRecords,
	) (domain.MedicalRecords, error)
}
