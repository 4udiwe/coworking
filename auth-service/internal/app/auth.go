package app

import (
	"crypto/rsa"

	"github.com/4udiwe/coworking/auth-service/internal/auth"
	"github.com/4udiwe/coworking/auth-service/pgk/jwt_validator"
	"github.com/sirupsen/logrus"
)

func (app *App) PrivateKey() *rsa.PrivateKey {
	if app.privateKey != nil {
		return app.privateKey
	}
	privateKey, err := auth.LoadPrivateKeyFromPEM(app.cfg.Auth.PrivateKey)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load private key")
	}
	app.privateKey = privateKey
	return app.privateKey
}

func (app *App) Auth() *auth.Auth {
	if app.auth != nil {
		return app.auth
	}
	app.auth = auth.New(
		app.PrivateKey(),
		app.cfg.App.Name,
		app.cfg.Auth.AccessTokenTTL,
		app.cfg.Auth.RefreshTokenTTL,
	)
	return app.auth
}

func (app *App) JwtValidator() *jwt_validator.Validator {
	if app.jwtValidator != nil {
		return app.jwtValidator
	}
	app.jwtValidator = jwt_validator.NewValidator(app.Auth().PublicKey)
	return app.jwtValidator
}
