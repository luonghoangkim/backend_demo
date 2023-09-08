package handler

import (
	"backend_myblog/banana"
	"backend_myblog/log"
	"backend_myblog/model"
	"backend_myblog/model/req"
	repository "backend_myblog/repository"
	"backend_myblog/security"
	"net/http"

	"github.com/golang-jwt/jwt"
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
 
	if err := c.Validate(req); err != nil {
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

	//gen token 
	token,err := security.GenToken(user)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user.Token = token
 
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       user,
	})
}

// HandleSignIn xử lý yêu cầu đăng nhập người dùng.
func (u *UserHandler) HandleSignIn(c echo.Context) error {
	req := req.ReqSignIn{}

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
	if err := c.Validate(req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user, err := u.UserRepo.CheckLogin(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// check pass 
	isTheSame := security.ComparePasswords(user.PassWord, []byte(req.Password))
	if !isTheSame {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Đăng nhập thất bại",
			Data:       nil,
		})
	}

	//gen token 
	token,err := security.GenToken(user)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user.Token = token

	// Trả về thông tin người dùng giả mạo để kiểm tra yêu cầu đăng nhập. 
	return c.JSON(http.StatusOK, model.Response{
			StatusCode: http.StatusOK,
			Message:    "Xử lý thành công",
			Data:       user,
		})
} 

func (u *UserHandler) Profile(c echo.Context) error {
	tokenData := c.Get("user").(*jwt.Token)
	clams := tokenData.Claims.(*model.JwtCustomClaim)

	user, err := u.UserRepo.SelectUserByID(c.Request().Context(), clams.UserID)

	if err != nil {
		if err == banana.UserNotFound {
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}

		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
		
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       user,
	})

}

func (u UserHandler) UpdateProfile(c echo.Context) error {
	req := req.ReqUpdateUser{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomClaim)
	user := model.User{
		UserID:   claims.UserID,
		FullName: req.FullName,
		Email:    req.Email,
	}

	user, err = u.UserRepo.UpdateUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, model.Response{
		StatusCode: http.StatusCreated,
		Message:    "Xử lý thành công",
		Data:       user,
	})
}