package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/internal/application/medical/service"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

type medicalHandler struct {
	medicalService service.MedicalServiceContract
}

func NewMedicalHandler(
	router fiber.Router,
	jwtMiddleware fiber.Handler,
	medicalService service.MedicalServiceContract,
) {
	handler := medicalHandler{
		medicalService: medicalService,
	}

	medicalRouter := router.Group("/medical", jwtMiddleware)
	medicalRouter.Post("/patient", handler.RecordPatient)
	medicalRouter.Get("/patient", handler.GetPatients)
	medicalRouter.Post("/record", handler.SaveMedicalRecord)
	medicalRouter.Get("/record", handler.GetMedicalRecords)
}

func (h medicalHandler) RecordPatient(c *fiber.Ctx) error {
	callerInfo := "[medicalHandler.RecordPatient]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	req := recordPatientReqAcquire()
	defer recordPatientReqRelease(req)

	if err := c.BodyParser(req); err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if err := req.validate(); err != nil {
		l.Error("error validating request", zap.Error(err))
		return errBadRequest{err: err}
	}

	patient := domain.PatientAcquire()
	defer domain.PatientRelease(patient)

	patient.ID = string(*req.IdentityNumber)
	patient.PhoneNumber = req.PhoneNumber
	patient.Name = req.Name
	patient.BirthDate = req.birthDate
	patient.Gender = req.Gender
	patient.ImgURL = req.ImgURL

	err := h.medicalService.RecordPatient(userCtx, patient)
	if err != nil {
		l.Error("failed to record patient", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "Patient recorded successfully"

	return c.Status(http.StatusCreated).JSON(res)
}

func (h medicalHandler) GetPatients(c *fiber.Ctx) error {
	callerInfo := "[medicalHandler.GetPatients]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	query := queryPatientAcquire()
	defer queryPatientRelease(query)

	if err := c.QueryParser(query); err != nil {
		l.Error("error parsing query params", zap.Error(err))
		return errBadRequest{err: err}
	}

	query.validate()

	filter := domain.FilterPatientAcquire()
	defer domain.FilterPatientRelease(filter)

	filter.ID = string(query.IdentityNumber)
	filter.Limit = query.Limit
	filter.Offset = query.Offset
	filter.Name = query.Name
	filter.PhoneNumber = query.PhoneNumber
	filter.CreatedAt = query.CreatedAt

	patients := domain.PatientsAcquire()
	defer domain.PatientsRelease(patients)

	patients, err := h.medicalService.GetPatients(userCtx, filter, patients)
	if err != nil {
		l.Error("failed to get patients", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "Patients retrieved successfully"

	patientsRes := getPatientsResAcquire()
	defer getPatientsResRelease(patientsRes)

	for _, patient := range patients {
		patientsRes = append(patientsRes, getPatientRes{
			IdentityNumber: idNumber(patient.ID),
			PhoneNumber:    "+" + patient.PhoneNumber,
			Name:           patient.Name,
			BirthDate:      patient.BirthDate.Format(dateFormat),
			Gender:         patient.Gender,
			CreatedAt:      patient.CreatedAt.Format(dateFormat),
		})
	}

	res.Data = patientsRes

	return c.JSON(res)
}

func (h medicalHandler) SaveMedicalRecord(c *fiber.Ctx) error {
	callerInfo := "[medicalHandler.SaveMedicalRecord]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	req := medicalRecordAcquire()
	defer medicalRecordRelease(req)

	if err := c.BodyParser(req); err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if err := req.validate(); err != nil {
		l.Error("error validating request", zap.Error(err))
		return errBadRequest{err: err}
	}

	record := domain.MedicalRecordAcquire()
	defer domain.MedicalRecordRelease(record)

	user := domain.UserAcquire()
	defer domain.UserRelease(user)
	*user = c.Locals(domain.UserFromToken).(domain.User)

	record.PatientID = string(*req.IdentityNumber)
	record.Symptoms = req.Symptoms
	record.Medications = req.Medications
	record.StaffID = user.ID
	record.StaffNIP = user.NIP
	record.StaffName = user.Name

	err := h.medicalService.SaveMedicalRecord(userCtx, record, user)
	if err != nil {
		l.Error("failed to save medical record", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "Medical record saved successfully"

	return c.Status(http.StatusCreated).JSON(res)
}

func (h medicalHandler) GetMedicalRecords(c *fiber.Ctx) error {
	callerInfo := "[medicalHandler.GetMedicalRecords]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	query := queryRecordAcquire()
	defer queryRecordRelease(query)

	query.PatientID = c.QueryInt("identityDetail.identityNumber", 0)
	query.StaffID = c.Query("createdBy.userId", "")
	query.StaffNIP = c.Query("createdBy.nip", "")
	query.Limit = c.QueryInt("limit", 0)
	query.Offset = c.QueryInt("offset", 0)
	query.CreatedAt = c.Query("createdAt", "")
	query.validate()

	filter := domain.FilterMedicalRecordAcquire()
	defer domain.FilterMedicalRecordRelease(filter)

	filter.PatientID = query.patientID
	filter.StaffID = query.staffID
	filter.StaffNIP = query.StaffNIP
	filter.Limit = query.Limit
	filter.Offset = query.Offset
	filter.CreatedAt = query.CreatedAt

	records := domain.MedicalRecordsAcquire()
	defer domain.MedicalRecordsRelease(records)

	records, err := h.medicalService.GetMedicalRecords(userCtx, filter, records)
	if err != nil {
		l.Error("failed to get medical records", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "Medical records retrieved successfully"

	recordsRes := getRecordsResAcquire()
	defer getRecordsResRelease(recordsRes)
	var nip int

	for _, record := range records {
		nip, _ = strconv.Atoi(record.StaffNIP)

		recordsRes = append(recordsRes, getRecordRes{
			IdentityDetail: identityDetail{
				IdentityNumber:      idNumber(record.PatientID),
				PhoneNumber:         "+" + record.PatientPhoneNumber,
				Name:                record.PatientName,
				BirthDate:           record.PatientBirthDate.Format(dateFormat),
				Gender:              record.PatientGender,
				IdentityCardScanImg: record.PatientImgURL,
			},
			Symptoms:    record.Symptoms,
			Medications: record.Medications,
			CreatedAt:   record.CreatedAt.Format(dateFormat),
			CreatedBy: createdBy{
				Nip:    uint(nip),
				Name:   record.StaffName,
				UserId: record.StaffID,
			},
		})
	}

	res.Data = recordsRes

	return c.JSON(res)
}
