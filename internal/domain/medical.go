package domain

import (
	"net/http"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	GenderMale   = "male"
	GenderFemale = "female"
)

var PatientPool = sync.Pool{
	New: func() any {
		return new(Patient)
	},
}

func PatientAcquire() *Patient {
	return PatientPool.Get().(*Patient)
}

func PatientRelease(t *Patient) {
	*t = Patient{}
	PatientPool.Put(t)
}

type Patient struct {
	ID          string
	PhoneNumber string
	Name        string
	BirthDate   time.Time
	Gender      string
	ImgURL      string
	CreatedAt   time.Time
}

const patientsInitCap = 5

var PatientsPool = sync.Pool{
	New: func() any {
		return make(Patients, 0, patientsInitCap)
	},
}

func PatientsAcquire() Patients {
	return PatientsPool.Get().(Patients)
}

func PatientsRelease(t Patients) {
	t = t[:0]
	PatientsPool.Put(t) // nolint:staticcheck
}

type Patients []Patient

var FilterPatientPool = sync.Pool{
	New: func() any {
		return new(FilterPatient)
	},
}

func FilterPatientAcquire() *FilterPatient {
	return FilterPatientPool.Get().(*FilterPatient)
}

func FilterPatientRelease(t *FilterPatient) {
	*t = FilterPatient{}
	FilterPatientPool.Put(t)
}

type FilterPatient struct {
	ID          string
	Limit       int
	Offset      int
	Name        string
	PhoneNumber string
	CreatedAt   string
}

var MedicalRecordPool = sync.Pool{
	New: func() any {
		return new(MedicalRecord)
	},
}

func MedicalRecordAcquire() *MedicalRecord {
	return MedicalRecordPool.Get().(*MedicalRecord)
}

func MedicalRecordRelease(t *MedicalRecord) {
	*t = MedicalRecord{}
	MedicalRecordPool.Put(t)
}

type MedicalRecord struct {
	ID                 ulid.ULID
	PatientID          string
	PatientPhoneNumber string
	PatientName        string
	PatientBirthDate   time.Time
	PatientGender      string
	PatientImgURL      string
	Symptoms           string
	Medications        string
	StaffID            ulid.ULID
	StaffNIP           string
	StaffName          string
	CreatedAt          time.Time
}

var FilterMedicalRecordPool = sync.Pool{
	New: func() any {
		return new(FilterMedicalRecord)
	},
}

func FilterMedicalRecordAcquire() *FilterMedicalRecord {
	return FilterMedicalRecordPool.Get().(*FilterMedicalRecord)
}

func FilterMedicalRecordRelease(t *FilterMedicalRecord) {
	*t = FilterMedicalRecord{}
	FilterMedicalRecordPool.Put(t)
}

type FilterMedicalRecord struct {
	PatientID string
	StaffID   ulid.ULID
	StaffNIP  string
	Limit     int
	Offset    int
	CreatedAt string
}

const medicalRecordsInitCap = 5

var MedicalRecordsPool = sync.Pool{
	New: func() any {
		return make(MedicalRecords, 0, medicalRecordsInitCap)
	},
}

func MedicalRecordsAcquire() MedicalRecords {
	return MedicalRecordsPool.Get().(MedicalRecords)
}

func MedicalRecordsRelease(t MedicalRecords) {
	t = t[:0]
	MedicalRecordsPool.Put(t) // nolint:staticcheck
}

type MedicalRecords []MedicalRecord

type ErrDuplicatePatient struct{}

func (e ErrDuplicatePatient) Error() string {
	return "Patient already registered"
}

func (e ErrDuplicatePatient) Status() int {
	return http.StatusConflict
}

type ErrPatientNotFound struct{}

func (e ErrPatientNotFound) Error() string {
	return "Patient not found"
}

func (e ErrPatientNotFound) Status() int {
	return http.StatusNotFound
}
