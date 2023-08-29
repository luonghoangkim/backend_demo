package handler

import (
	"backend_myblog/log"
	"backend_myblog/model"
	"backend_myblog/model/req"
	repository "backend_myblog/repository"
	"backend_myblog/security"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	uuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserRepo repository.UserRepo
}
 
// HandleSignUp xử lý yêu cầu đăng ký người dùng.
func (u *UserHandler) HandleSignUp(c echo.Context) error {
	// Tạo một biến req để chứa dữ liệu từ yêu cầu đăng ký.
	req := req.ReqSignUp{}

	// Kiểm tra và liên kết dữ liệu từ yêu cầu đến biến req.
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// Tạo một validator mới để kiểm tra dữ liệu.
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// Băm và mã hóa mật khẩu người dùng.
	hash := security.HashAndSalt([]byte(req.Password))
	role := model.MENBER.String()

	// Tạo một User ID mới sử dụng UUID.
	userId, err := uuid.NewUUID()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// Tạo một biến User và lưu thông tin người dùng vào cơ sở dữ liệu.
	user := model.User{
		UserID:   userId.String(),
		FullName: req.FullName,
		Email:    req.Email,
		PassWord: hash,
		Role:     role,
		Token:    "",
	}

	user, err = u.UserRepo.SaveUser(c.Request().Context(), user)
	
	if err != nil {
		return c.JSON(http.StatusConflict, model.Response{
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user.PassWord = "";
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       user,
	})
}

// HandleSignIn xử lý yêu cầu đăng nhập người dùng.
func (u *UserHandler) HandleSignIn(c echo.Context) error {
	// Trả về thông tin người dùng giả mạo để kiểm tra yêu cầu đăng nhập.
	return c.JSON(http.StatusOK, echo.Map{
		"user":  "Lương",
		"email": "hoangkimluong192@gmail.com",
	})
} 
