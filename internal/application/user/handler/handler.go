package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"github.com/j03hanafi/halo-suster/common/logger"
	"github.com/j03hanafi/halo-suster/internal/application/user/service"
	"github.com/j03hanafi/halo-suster/internal/domain"
)

const userIDFromParam = "id"

type userHandler struct {
	userService service.UserServiceContract
}

func NewUserHandler(router fiber.Router, jwtMiddleware fiber.Handler, userService service.UserServiceContract) {
	handler := userHandler{userService: userService}

	authRouter := router.Group("/user")
	authRouter.Post("/it/register", handler.RegisterIT)
	authRouter.Post("/it/login", handler.LoginIT)
	authRouter.Post("/nurse/login", handler.LoginNurse)
	authRouter.Get("", jwtMiddleware, itStaffAccess, handler.GetUsers)

	nurseRouter := router.Group("/user/nurse", jwtMiddleware, itStaffAccess)
	nurseRouter.Post("/register", handler.RegisterNurse)
	nurseRouter.Put("/:"+userIDFromParam, handler.UpdateNurse)
	nurseRouter.Delete("/:"+userIDFromParam, handler.DeleteNurse)
	nurseRouter.Post("/:"+userIDFromParam+"/access", handler.UpdateAccess)
}

func (h userHandler) RegisterIT(c *fiber.Ctx) error {
	callerInfo := "[userHandler.RegisterIT]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	req := registerITReqAcquire()
	defer registerITReqRelease(req)

	if err := c.BodyParser(req); err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if err := req.validate(nipITFirstDigit); err != nil {
		l.Error("error validating request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	user := domain.UserAcquire()
	defer domain.UserRelease(user)

	user.NIP = string(*req.NIP)
	user.Name = req.Name
	user.Password = req.Password

	err := h.userService.RegisterIT(userCtx, user)
	if err != nil {
		l.Error("error registering IT user", zap.Error(err))
		return err
	}

	token, err := h.userService.GenerateToken(userCtx, user)
	if err != nil {
		l.Error("error generating token", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "User registered successfully"

	data := authUserResAcquire()
	defer authUserResRelease(data)

	data.UserID = user.ID
	data.Name = user.Name
	data.NIP = nip(user.NIP)
	data.AccessToken = token
	res.Data = data

	return c.Status(http.StatusCreated).JSON(res)
}

func (h userHandler) LoginIT(c *fiber.Ctx) error {
	callerInfo := "[userHandler.LoginIT]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	req := loginReqAcquire()
	defer loginReqRelease(req)

	if err := c.BodyParser(req); err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if err := req.validate(nipITFirstDigit); err != nil {
		l.Error("error validating request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if !strings.HasPrefix(string(*req.NIP), nipITFirstDigit) {
		return new(domain.ErrInvalidNIP)
	}

	user := domain.UserAcquire()
	defer domain.UserRelease(user)

	user.NIP = string(*req.NIP)
	user.Password = req.Password

	user, err := h.userService.LoginIT(userCtx, user)
	if err != nil {
		l.Error("error logging in IT user", zap.Error(err))
		return err
	}

	token, err := h.userService.GenerateToken(userCtx, user)
	if err != nil {
		l.Error("error generating token", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "User logged in successfully"

	data := authUserResAcquire()
	defer authUserResRelease(data)

	data.UserID = user.ID
	data.Name = user.Name
	data.NIP = nip(user.NIP)
	data.AccessToken = token
	res.Data = data

	return c.JSON(res)
}

func (h userHandler) LoginNurse(c *fiber.Ctx) error {
	callerInfo := "[userHandler.LoginNurse]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	req := loginReqAcquire()
	defer loginReqRelease(req)

	if err := c.BodyParser(req); err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if err := req.validate(nipNurseFirstDigit); err != nil {
		l.Error("error validating request body", zap.Error(err))

		return errBadRequest{err: err}
	}

	if !strings.HasPrefix(string(*req.NIP), nipNurseFirstDigit) {
		return new(domain.ErrInvalidNIP)
	}

	user := domain.UserAcquire()
	defer domain.UserRelease(user)

	user.NIP = string(*req.NIP)
	user.Password = req.Password

	_, err := h.userService.LoginNurse(userCtx, user)
	if err != nil {
		l.Error("error logging in nurse user", zap.Error(err))
		return err
	}

	token, err := h.userService.GenerateToken(userCtx, user)
	if err != nil {
		l.Error("error generating token", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "User logged in successfully"

	data := authUserResAcquire()
	defer authUserResRelease(data)

	data.UserID = user.ID
	data.Name = user.Name
	data.NIP = nip(user.NIP)
	data.AccessToken = token
	res.Data = data

	return c.JSON(res)
}

func (h userHandler) RegisterNurse(c *fiber.Ctx) error {
	callerInfo := "[userHandler.RegisterNurse]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	req := registerNurseReqAcquire()
	defer registerNurseReqRelease(req)

	if err := c.BodyParser(req); err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if err := req.validate(nipNurseFirstDigit); err != nil {
		l.Error("error validating request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	user := domain.UserAcquire()
	defer domain.UserRelease(user)

	user.NIP = string(*req.NIP)
	user.Name = req.Name
	user.ImgURL = req.ImgURL

	err := h.userService.RegisterNurse(userCtx, user)
	if err != nil {
		l.Error("error registering nurse user", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "User registered successfully"

	data := authUserResAcquire()
	defer authUserResRelease(data)

	data.UserID = user.ID
	data.Name = user.Name
	data.NIP = nip(user.NIP)
	res.Data = data

	return c.Status(http.StatusCreated).JSON(res)
}

func (h userHandler) UpdateNurse(c *fiber.Ctx) error {
	callerInfo := "[userHandler.UpdateNurse]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userID, err := ulid.Parse(c.Params(userIDFromParam))
	if err != nil {
		l.Error("error parsing userIDParam", zap.Error(err))
		return errBadRequest{err: err}
	}

	req := updateNurseReqAcquire()
	defer updateNurseReqRelease(req)

	if err = c.BodyParser(req); err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if err = req.validate(nipNurseFirstDigit); err != nil {
		l.Error("error validating request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if !strings.HasPrefix(string(*req.NIP), nipNurseFirstDigit) {
		return new(domain.ErrInvalidNIP)
	}

	user := domain.UserAcquire()
	defer domain.UserRelease(user)

	user.ID = userID
	user.NIP = string(*req.NIP)
	user.Name = req.Name

	err = h.userService.UpdateNurse(userCtx, user)
	if err != nil {
		l.Error("error updating nurse user", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "User updated successfully"

	return c.JSON(res)
}

func (h userHandler) DeleteNurse(c *fiber.Ctx) error {
	callerInfo := "[userHandler.DeleteNurse]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userID, err := ulid.Parse(c.Params(userIDFromParam))
	if err != nil {
		l.Error("error parsing userIDParam", zap.Error(err))
		return errBadRequest{err: err}
	}

	user := domain.UserAcquire()
	defer domain.UserRelease(user)

	user.ID = userID

	err = h.userService.DeleteNurse(userCtx, user)
	if err != nil {
		l.Error("error deleting nurse user", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "User deleted successfully"

	return c.JSON(res)
}

func (h userHandler) UpdateAccess(c *fiber.Ctx) error {
	callerInfo := "[userHandler.UpdateAccess]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userID, err := ulid.Parse(c.Params(userIDFromParam))
	if err != nil {
		l.Error("error parsing userIDParam", zap.Error(err))
		return errBadRequest{err: err}
	}

	req := updateAccessReqAcquire()
	defer updateAccessReqRelease(req)

	if err = c.BodyParser(req); err != nil {
		l.Error("error parsing request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	if err = req.validate(); err != nil {
		l.Error("error validating request body", zap.Error(err))
		return errBadRequest{err: err}
	}

	user := domain.UserAcquire()
	defer domain.UserRelease(user)

	user.ID = userID
	user.Password = req.Password

	err = h.userService.UpdateAccess(userCtx, user)
	if err != nil {
		l.Error("error deleting nurse user", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "Access updated successfully"

	return c.JSON(res)
}

func (h userHandler) GetUsers(c *fiber.Ctx) error {
	callerInfo := "[userHandler.GetUsers]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	query := queryParamAcquire()
	defer queryParamRelease(query)

	if err := c.QueryParser(query); err != nil {
		l.Error("error parsing query params", zap.Error(err))
		return errBadRequest{err: err}
	}

	query.validate()

	filter := domain.FilterUserAcquire()
	defer domain.FilterUserRelease(filter)

	filter.UserID = query.uid
	filter.Limit = int(query.Limit)
	filter.Offset = int(query.Offset)
	filter.Name = query.Name
	filter.NIP = query.nip
	filter.Role = query.Role
	filter.CreatedAt = query.CreatedAt

	users := domain.UsersAcquire()
	defer domain.UsersRelease(users)

	users, err := h.userService.GetUsers(userCtx, filter, users)
	if err != nil {
		l.Error("error getting users", zap.Error(err))
		return err
	}

	res := baseResponseAcquire()
	defer baseResponseRelease(res)

	res.Message = "Users retrieved successfully"

	usersRes := getUsersResAcquire()
	defer getUsersResRelease(usersRes)

	userRes := getUserResAcquire()
	defer getUserResRelease(userRes)

	for _, user := range users {
		userRes.UserID = user.ID
		userRes.NIP = nip(user.NIP)
		userRes.Name = user.Name
		userRes.CreatedAt = user.CreatedAt.Format(time.DateOnly)

		usersRes = append(usersRes, *userRes)
	}

	res.Data = usersRes

	return c.JSON(res)
}

func itStaffAccess(c *fiber.Ctx) error {
	user := domain.UserAcquire()
	defer domain.UserRelease(user)
	*user = c.Locals(domain.UserFromToken).(domain.User)

	if user.Role != domain.RoleIT {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
		})
	}

	return c.Next()
}
