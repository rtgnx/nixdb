package main

import (
	"fmt"
	"net/http"

	"github.com/Reverse-Labs/nixdb"
	"github.com/labstack/echo"
)

type HTTP struct {
	db nixdb.Database
}

func (h HTTP) GETLogin(ctx echo.Context) error {

	if user, ok := ctx.Get("user").(nixdb.PasswdEntry); ok {
		AttachProxyAuthHeaders(ctx, user)
		return ctx.JSON(http.StatusOK, user)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h HTTP) GETUsers(ctx echo.Context) error {

	if err := h.db.Update(); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, h.db.Users)
}

func (h HTTP) GETGroups(ctx echo.Context) error {

	if err := h.db.Update(); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, h.db.Groups)
}

func (h HTTP) GETHosts(ctx echo.Context) error {
	if err := h.db.Update(); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, h.db.Hosts)
}

func AttachProxyAuthHeaders(ctx echo.Context, user nixdb.PasswdEntry) {
	ctx.Response().Header().Set("X-Forwarded-User", user.Name)
	ctx.Response().Header().Add("X-Forwarded-FullName", user.Fullname)
	ctx.Response().Header().Add("X-Forwarded-Uid", fmt.Sprintf("%d", user.UID))
	ctx.Response().Header().Add("X-Forwarded-Gid", fmt.Sprintf("%d", user.GID))
}
