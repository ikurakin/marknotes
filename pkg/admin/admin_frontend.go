package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/marknotes/pkg/models"
)

type AdminFrontend struct {
	repo models.PostRepository
}

func NewAdminFrontend(g *echo.Group, repo models.PostRepository) {
	fe := &AdminFrontend{
		repo: repo,
	}

	g.GET("", fe.Index)
}

func (fe *AdminFrontend) Index(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/posts_index")
}
