package handler

import (
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Repository repository.RepositoryInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
	}
}

func sendErrorResponse(ctx echo.Context, httpCode int, err error) error {
	var errorResp generated.ErrorResponse

	errorResp.Message = err.Error()
	return ctx.JSON(httpCode, errorResp)

}
