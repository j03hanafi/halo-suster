package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/oklog/ulid/v2"
	"go.uber.org/multierr"

	"github.com/j03hanafi/halo-suster/internal/domain"
)

type errBadRequest struct {
	err error
}

func (e errBadRequest) Error() string {
	return e.err.Error()
}

func (e errBadRequest) Status() int {
	return http.StatusBadRequest
}

var baseResponsePool = sync.Pool{
	New: func() any {
		return new(baseResponse)
	},
}

func baseResponseAcquire() *baseResponse {
	return baseResponsePool.Get().(*baseResponse)
}

func baseResponseRelease(t *baseResponse) {
	*t = baseResponse{}
	baseResponsePool.Put(t)
}

type baseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type idNumber string

func (n *idNumber) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.New("identityNumber is required")
	}

	var jsonID int
	if err := json.Unmarshal(b, &jsonID); err != nil {
		return errors.New("identityNumber must be a number")
	}
	*n = idNumber(strconv.Itoa(jsonID))
	return nil
}

func (n *idNumber) MarshalJSON() ([]byte, error) {
	jsonID, err := strconv.Atoi(string(*n))
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonID)
}

func (n *idNumber) validate() error {
	var errs error
	const idNumberLength = 16

	if len(*n) != idNumberLength {
		errs = multierr.Append(errs, errors.New("identityNumber must have 16 characters"))
		return errs
	}

	return nil
}

var recordPatientReqPool = sync.Pool{
	New: func() any {
		return new(recordPatientReq)
	},
}

func recordPatientReqAcquire() *recordPatientReq {
	return recordPatientReqPool.Get().(*recordPatientReq)
}

func recordPatientReqRelease(t *recordPatientReq) {
	*t = recordPatientReq{}
	recordPatientReqPool.Put(t)
}

type recordPatientReq struct {
	IdentityNumber *idNumber `json:"identityNumber"`
	PhoneNumber    string    `json:"phoneNumber"`
	Name           string    `json:"name"`
	BirthDate      string    `json:"birthDate"`
	birthDate      time.Time
	Gender         string `json:"gender"`
	ImgURL         string `json:"identityCardScanImg"`
}

func (r *recordPatientReq) validate() error {
	var errs error

	if r.IdentityNumber == nil {
		errs = multierr.Append(errs, errors.New("identityNumber is required"))
	} else {
		errs = multierr.Append(errs, r.IdentityNumber.validate())
	}

	if r.PhoneNumber == "" {
		errs = multierr.Append(errs, errors.New("phoneNumber is required"))
	} else if len(r.PhoneNumber) < 10 || len(r.PhoneNumber) > 15 {
		errs = multierr.Append(errs, errors.New("phoneNumber must have 10 to 15 characters"))
	} else if !strings.HasPrefix(r.PhoneNumber, "+62") {
		errs = multierr.Append(errs, errors.New("phoneNumber must start with +62"))
	} else {
		r.PhoneNumber = strings.TrimPrefix(r.PhoneNumber, "+")
	}

	if r.Name == "" {
		errs = multierr.Append(errs, errors.New("name is required"))
	} else if len(r.Name) < 3 || len(r.Name) > 30 {
		errs = multierr.Append(errs, errors.New("name must have 3 to 30 characters"))
	}

	if r.BirthDate == "" {
		errs = multierr.Append(errs, errors.New("birthDate is required"))
	} else {
		if birthDate, err := time.Parse("2006-01-02T15:04:05.999Z", r.BirthDate); err != nil {
			errs = multierr.Append(errs, errors.New("birthDate must be in format yyyy-mm-dd"))
		} else {
			r.birthDate = birthDate
		}
	}

	if r.Gender == "" {
		errs = multierr.Append(errs, errors.New("gender is required"))
	} else if r.Gender != domain.GenderMale && r.Gender != domain.GenderFemale {
		errs = multierr.Append(errs, errors.New("gender must be either male or female"))
	}

	if r.ImgURL == "" {
		errs = multierr.Append(errs, errors.New("identity card scan image URL is required"))
	} else if !govalidator.IsURL(r.ImgURL) {
		errs = multierr.Append(errs, errors.New("identity card scan image URL must be a valid URL"))
	} else {
		u, err := url.Parse(r.ImgURL)
		if err != nil {
			errs = multierr.Append(errs, errors.New("identity card scan image URL must be a valid URL"))
		}
		if u != nil && !strings.Contains(u.Host, ".") {
			errs = multierr.Append(errs, errors.New("identity card scan image URL must be a valid URL"))
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}

var queryPatientPool = sync.Pool{
	New: func() any {
		return new(queryPatient)
	},
}

func queryPatientAcquire() *queryPatient {
	return queryPatientPool.Get().(*queryPatient)
}

func queryPatientRelease(t *queryPatient) {
	*t = queryPatient{}
	queryPatientPool.Put(t)
}

type queryPatient struct {
	IdentityNumber idNumber `query:"identityNumber"`
	Limit          int      `query:"limit"`
	Offset         int      `query:"offset"`
	Name           string   `query:"name"`
	PhoneNumber    string   `query:"phoneNumber"`
	CreatedAt      string   `query:"createdAt"`
}

func (r *queryPatient) validate() {
	if r.CreatedAt != "" && r.CreatedAt != "asc" && r.CreatedAt != "desc" {
		r.CreatedAt = ""
	}
}

type getPatientRes struct {
	IdentityNumber idNumber `json:"identityNumber"`
	PhoneNumber    string   `json:"phoneNumber"`
	Name           string   `json:"name"`
	BirthDate      string   `json:"birthDate"`
	Gender         string   `json:"gender"`
	CreatedAt      string   `json:"createdAt"`
}

const patientsInitCap = 5

var getPatientsResPool = sync.Pool{
	New: func() any {
		return make(getPatientsRes, 0, patientsInitCap)
	},
}

func getPatientsResAcquire() getPatientsRes {
	return getPatientsResPool.Get().(getPatientsRes)
}

func getPatientsResRelease(t getPatientsRes) {
	t = t[:0]
	getPatientsResPool.Put(t) // nolint:staticcheck
}

type getPatientsRes []getPatientRes

var medicalRecordPool = sync.Pool{
	New: func() any {
		return new(medicalRecord)
	},
}

func medicalRecordAcquire() *medicalRecord {
	return medicalRecordPool.Get().(*medicalRecord)
}

func medicalRecordRelease(t *medicalRecord) {
	*t = medicalRecord{}
	medicalRecordPool.Put(t)
}

type medicalRecord struct {
	IdentityNumber *idNumber `json:"identityNumber"`
	Symptoms       string    `json:"symptoms"`
	Medications    string    `json:"medications"`
}

func (r medicalRecord) validate() error {
	var errs error

	if r.IdentityNumber == nil {
		errs = multierr.Append(errs, errors.New("identityNumber is required"))
	} else {
		errs = multierr.Append(errs, r.IdentityNumber.validate())
	}

	if r.Symptoms == "" {
		errs = multierr.Append(errs, errors.New("symptoms is required"))
	} else if len(r.Symptoms) < 1 || len(r.Symptoms) > 2000 {
		errs = multierr.Append(errs, errors.New("symptoms must have 1 to 2000 characters"))
	}

	if r.Medications == "" {
		errs = multierr.Append(errs, errors.New("medications is required"))
	} else if len(r.Medications) < 1 || len(r.Medications) > 2000 {
		errs = multierr.Append(errs, errors.New("medications must have 1 to 2000 characters"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

var queryRecordPool = sync.Pool{
	New: func() any {
		return new(queryRecord)
	},
}

func queryRecordAcquire() *queryRecord {
	return queryRecordPool.Get().(*queryRecord)
}

func queryRecordRelease(t *queryRecord) {
	*t = queryRecord{}
	queryRecordPool.Put(t)
}

type queryRecord struct {
	PatientID int `query:"identityDetail.identityNumber"`
	patientID string
	StaffID   string `query:"createdBy.userId"`
	staffID   ulid.ULID
	StaffNIP  string `query:"createdBy.nip"`
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
	CreatedAt string `query:"createdAt"`
}

func (r *queryRecord) validate() {
	if r.PatientID != 0 {
		r.patientID = strconv.Itoa(r.PatientID)
	}

	if r.StaffID != "" {
		r.staffID, _ = ulid.Parse(r.StaffID)
	}

	if r.CreatedAt != "" && r.CreatedAt != "asc" && r.CreatedAt != "desc" {
		r.CreatedAt = ""
	}
}

type identityDetail struct {
	IdentityNumber      idNumber `json:"identityNumber"`
	PhoneNumber         string   `json:"phoneNumber"`
	Name                string   `json:"name"`
	BirthDate           string   `json:"birthDate"`
	Gender              string   `json:"gender"`
	IdentityCardScanImg string   `json:"identityCardScanImg"`
}

type createdBy struct {
	Nip    uint      `json:"nip"`
	Name   string    `json:"name"`
	UserId ulid.ULID `json:"userId"`
}

type getRecordRes struct {
	IdentityDetail identityDetail `json:"identityDetail:"`
	Symptoms       string         `json:"symptoms"`
	Medications    string         `json:"medications"`
	CreatedAt      string         `json:"createdAt"`
	CreatedBy      createdBy      `json:"createdBy"`
}

const recordsInitCap = 5

var getRecordsResPool = sync.Pool{
	New: func() any {
		return make(getRecordsRes, 0, recordsInitCap)
	},
}

func getRecordsResAcquire() getRecordsRes {
	return getRecordsResPool.Get().(getRecordsRes)
}

func getRecordsResRelease(t getRecordsRes) {
	t = t[:0]
	getRecordsResPool.Put(t) // nolint:staticcheck
}

type getRecordsRes []getRecordRes
