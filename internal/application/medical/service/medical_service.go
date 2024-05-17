package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/internal/application/medical/repository"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

type MedicalService struct {
	medicalRepository repository.MedicalRepositoryContract
	contextTimeout    time.Duration
}

func NewMedicalService(timeout time.Duration, medicalRepository repository.MedicalRepositoryContract) *MedicalService {
	return &MedicalService{
		medicalRepository: medicalRepository,
		contextTimeout:    timeout,
	}
}

func (s MedicalService) RecordPatient(ctx context.Context, patient *domain.Patient) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[MedicalService.RecordPatient]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	err := s.medicalRepository.RecordPatient(ctx, patient)
	if err != nil {
		l.Error("failed to record patient", zap.Error(err))
		return err
	}

	return nil
}

func (s MedicalService) GetPatients(
	ctx context.Context,
	filter *domain.FilterPatient,
	patients domain.Patients,
) (domain.Patients, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[MedicalService.GetPatients]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	patients, err := s.medicalRepository.GetPatients(ctx, filter, patients)
	if err != nil {
		l.Error("failed to get patients", zap.Error(err))
		return nil, err
	}

	return patients, nil
}

func (s MedicalService) SaveMedicalRecord(ctx context.Context, record *domain.MedicalRecord, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[MedicalService.SaveMedicalRecord]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	record.StaffID = user.ID
	record.StaffNIP = user.NIP

	err := s.medicalRepository.SaveMedicalRecord(ctx, record)
	if err != nil {
		l.Error("failed to save medical record", zap.Error(err))
		return err
	}

	return nil
}

func (s MedicalService) GetMedicalRecords(
	ctx context.Context,
	filter *domain.FilterMedicalRecord,
	records domain.MedicalRecords,
) (domain.MedicalRecords, error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	callerInfo := "[MedicalService.GetMedicalRecords]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	records, err := s.medicalRepository.GetMedicalRecords(ctx, filter, records)
	if err != nil {
		l.Error("failed to get medical records", zap.Error(err))
		return nil, err
	}

	return records, nil
}

var _ MedicalServiceContract = (*MedicalService)(nil)
