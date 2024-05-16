package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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

const (
	nipITFirstDigit    = "615"
	nipNurseFirstDigit = "303"
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

type nip string

func (n *nip) UnmarshalJSON(b []byte) error {
	var jsonNIP int
	if err := json.Unmarshal(b, &jsonNIP); err != nil {
		return errors.New("nip must be a number")
	}
	*n = nip(strconv.Itoa(jsonNIP))
	return nil
}

func (n *nip) MarshalJSON() ([]byte, error) {
	jsonNIP, err := strconv.Atoi(string(*n))
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonNIP)
}

func (n *nip) validate(firstDigit string) error {
	var errs error
	const nipLength = 13

	if len(*n) != nipLength {
		errs = multierr.Append(errs, errors.New("nip must have 13 characters"))
		return errs
	}

	if string(*n)[:3] != firstDigit {
		errs = multierr.Append(errs, fmt.Errorf("nip must have %s in the first three characters", firstDigit))
	}

	if string(*n)[3:4] != "1" && string(*n)[3:4] != "2" {
		errs = multierr.Append(errs, errors.New("nip must have 1 or 2 in the fourth character"))
	}

	year, err := strconv.Atoi(string(*n)[4:8])
	if err != nil || year < 2000 || year > time.Now().Year() {
		errs = multierr.Append(errs, errors.New("nip must have valid year from 2000 to current year"))
	}

	month, err := strconv.Atoi(string(*n)[8:10])
	if err != nil || month < 1 || month > 12 {
		errs = multierr.Append(errs, errors.New("nip must have valid month from 1 to 12"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

var registerITReqPool = sync.Pool{
	New: func() any {
		return new(registerITReq)
	},
}

func registerITReqAcquire() *registerITReq {
	return registerITReqPool.Get().(*registerITReq)
}

func registerITReqRelease(t *registerITReq) {
	*t = registerITReq{}
	registerITReqPool.Put(t)
}

type registerITReq struct {
	NIP      *nip   `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r registerITReq) validate(nipCheck string) error {
	var errs error

	if r.NIP == nil {
		errs = multierr.Append(errs, errors.New("nip is required"))
	} else {
		err := r.NIP.validate(nipCheck)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if r.Name == "" {
		errs = multierr.Append(errs, errors.New("name is required"))
	} else if len(r.Name) < 5 || len(r.Name) > 50 {
		errs = multierr.Append(errs, errors.New("name must have 5 to 50 characters"))
	}

	if r.Password == "" {
		errs = multierr.Append(errs, errors.New("password is required"))
	} else if len(r.Password) < 5 || len(r.Password) > 33 {
		errs = multierr.Append(errs, errors.New("password must have 5 to 33 characters"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

var authUserResPool = sync.Pool{
	New: func() any {
		return new(authUserRes)
	},
}

func authUserResAcquire() *authUserRes {
	return authUserResPool.Get().(*authUserRes)
}

func authUserResRelease(t *authUserRes) {
	*t = authUserRes{}
	authUserResPool.Put(t)
}

type authUserRes struct {
	UserID      ulid.ULID `json:"userId"`
	NIP         nip       `json:"nip"`
	Name        string    `json:"name"`
	AccessToken string    `json:"accessToken,omitempty"`
}

var loginReqPool = sync.Pool{
	New: func() any {
		return new(loginReq)
	},
}

func loginReqAcquire() *loginReq {
	return loginReqPool.Get().(*loginReq)
}

func loginReqRelease(t *loginReq) {
	*t = loginReq{}
	loginReqPool.Put(t)
}

type loginReq struct {
	NIP      *nip   `json:"nip"`
	Password string `json:"password"`
}

func (r loginReq) validate(nipCheck string) error {
	var errs error

	if r.NIP == nil {
		errs = multierr.Append(errs, errors.New("nip is required"))
	} else {
		err := r.NIP.validate(nipCheck)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if r.Password == "" {
		errs = multierr.Append(errs, errors.New("password is required"))
	} else if len(r.Password) < 5 || len(r.Password) > 33 {
		errs = multierr.Append(errs, errors.New("password must have 5 to 33 characters"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

var registerNurseReqPool = sync.Pool{
	New: func() any {
		return new(registerNurseReq)
	},
}

func registerNurseReqAcquire() *registerNurseReq {
	return registerNurseReqPool.Get().(*registerNurseReq)
}

func registerNurseReqRelease(t *registerNurseReq) {
	*t = registerNurseReq{}
	registerNurseReqPool.Put(t)
}

type registerNurseReq struct {
	NIP    *nip   `json:"nip"`
	Name   string `json:"name"`
	ImgURL string `json:"identityCardScanImg"`
}

func (r registerNurseReq) validate(nipCheck string) error {
	var errs error

	if r.NIP == nil {
		errs = multierr.Append(errs, errors.New("nip is required"))
	} else {
		err := r.NIP.validate(nipCheck)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if r.Name == "" {
		errs = multierr.Append(errs, errors.New("name is required"))
	} else if len(r.Name) < 5 || len(r.Name) > 50 {
		errs = multierr.Append(errs, errors.New("name must have 5 to 50 characters"))
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

var updateNurseReqPool = sync.Pool{
	New: func() any {
		return new(updateNurseReq)
	},
}

func updateNurseReqAcquire() *updateNurseReq {
	return updateNurseReqPool.Get().(*updateNurseReq)
}

func updateNurseReqRelease(t *updateNurseReq) {
	*t = updateNurseReq{}
	updateNurseReqPool.Put(t)
}

type updateNurseReq struct {
	NIP  *nip   `json:"nip"`
	Name string `json:"name"`
}

func (r updateNurseReq) validate(nipCheck string) error {
	var errs error

	if r.NIP == nil {
		errs = multierr.Append(errs, errors.New("nip is required"))
	} else {
		err := r.NIP.validate(nipCheck)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	if r.Name == "" {
		errs = multierr.Append(errs, errors.New("name is required"))
	} else if len(r.Name) < 5 || len(r.Name) > 50 {
		errs = multierr.Append(errs, errors.New("name must have 5 to 50 characters"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

var updateAccessReqPool = sync.Pool{
	New: func() any {
		return new(updateAccessReq)
	},
}

func updateAccessReqAcquire() *updateAccessReq {
	return updateAccessReqPool.Get().(*updateAccessReq)
}

func updateAccessReqRelease(t *updateAccessReq) {
	*t = updateAccessReq{}
	updateAccessReqPool.Put(t)
}

type updateAccessReq struct {
	Password string `json:"password"`
}

func (r updateAccessReq) validate() error {
	var errs error

	if r.Password == "" {
		errs = multierr.Append(errs, errors.New("password is required"))
	} else if len(r.Password) < 5 || len(r.Password) > 33 {
		errs = multierr.Append(errs, errors.New("password must have 5 to 33 characters"))
	}

	if errs != nil {
		return errs
	}

	return nil
}

var queryParamPool = sync.Pool{
	New: func() any {
		return new(queryParam)
	},
}

func queryParamAcquire() *queryParam {
	return queryParamPool.Get().(*queryParam)
}

func queryParamRelease(t *queryParam) {
	*t = queryParam{}
	queryParamPool.Put(t)
}

type queryParam struct {
	UserID    string `query:"userId"`
	uid       ulid.ULID
	Limit     uint   `query:"limit"`
	Offset    uint   `query:"offset"`
	Name      string `query:"name"`
	NIP       uint   `query:"nip"`
	nip       string
	Role      string `query:"role"`
	CreatedAt string `query:"createdAt"`
}

func (r *queryParam) validate() {
	if r.UserID != "" {
		r.uid, _ = ulid.Parse(r.UserID)
	}

	if r.NIP != 0 {
		r.nip = strconv.Itoa(int(r.NIP))
	}

	if r.Role != "" && r.Role != domain.RoleIT && r.Role != domain.RoleNurse {
		r.Role = ""
	}

	if r.CreatedAt != "" && r.CreatedAt != "asc" && r.CreatedAt != "desc" {
		r.CreatedAt = ""
	}
}

type getUserRes struct {
	UserID    ulid.ULID `json:"userId"`
	NIP       string    `json:"nip"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"createdAt"`
}

const usersInitCap = 5

var getUsersResPool = sync.Pool{
	New: func() any {
		return make(getUsersRes, 0, usersInitCap)
	},
}

func getUsersResAcquire() getUsersRes {
	return getUsersResPool.Get().(getUsersRes)
}

func getUsersResRelease(t getUsersRes) {
	t = t[:0]
	getUsersResPool.Put(t) // nolint:staticcheck
}

type getUsersRes []getUserRes
