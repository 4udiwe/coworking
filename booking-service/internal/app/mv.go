package app

import (
	"github.com/4udiwe/cowoking/booking-service/internal/api/middleware"
	"github.com/4udiwe/coworking/auth-service/pkg/jwt_validator"
	"github.com/sirupsen/logrus"
)

func (app *App) AuthMiddleware() *middleware.AuthMiddleware {
	if app.authMW != nil {
		return app.authMW
	}
	app.authMW = middleware.New(app.JwtValidator())
	return app.authMW
}

func (app *App) JwtValidator() *jwt_validator.Validator {
	if app.jwtValidator != nil {
		return app.jwtValidator
	}
	publicKey, err := jwt_validator.LoadPublicKey(app.cfg.Auth.PublicKey)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load public key")
	}
	app.jwtValidator = jwt_validator.NewValidator(publicKey)
	return app.jwtValidator
}
