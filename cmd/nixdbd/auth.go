package main

import (
	"fmt"
	"strings"

	"github.com/Reverse-Labs/nixdb"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func AuthSkipper(ctx echo.Context) bool {
	authHeader := ctx.Request().Header.Get(echo.HeaderAuthorization)
	return !strings.Contains(authHeader, "Basic") && len(authHeader) > 0
}

func AuthMiddleware(secret JWTSecretKey, db nixdb.Database, authorizedGroups []string) echo.MiddlewareFunc {
	return middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper:   AuthSkipper,
		Validator: BasicAuthValidator(secret, db, authorizedGroups),
		Realm:     middleware.DefaultBasicAuthConfig.Realm,
	})
}

func BasicAuthValidator(secret JWTSecretKey, db nixdb.Database, authorizedGroups []string) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {
		user, ok := db.Users.FindByName(username)

		if !ok {
			return false, fmt.Errorf("user: %s not found", username)
		}

		if !(nixdb.PAMAuth(username, password) && isAuthorized(db, username, authorizedGroups)) {
			return false, nil
		}

		if err := db.Update("passwd"); err != nil {
			return false, err
		}

		token, err := JWTToken(secret, NewJWTClaim(user, true, false))

		if err != nil {
			return false, fmt.Errorf("unable to sign token: %s", err.Error())
		}

		c.Response().Header().Add("Authorization", "Bearer "+token)
		c.Request().Header.Add("Authorization", "Bearer "+token)

		c.Set("user", user)

		return true, nil
	}
}

func isAuthorized(db nixdb.Database, username string, groups []string) bool {
	for _, authorizedGroup := range groups {
		if group, ok := db.Groups.FindByName(authorizedGroup); ok {
			if _, ok = group.Members(db.Users).FindByName(username); ok {
				return true
			}
		}
	}

	return false
}
