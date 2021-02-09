package firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/techartificer/swiftex/lib/errors"
	"google.golang.org/api/option"
)

var authClient *auth.Client

func Initialize() error {
	opt := option.WithCredentialsFile("swiftex-firebase.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}
	auth, err := app.Auth(context.Background())
	authClient = auth
	return err
}

func AuthClient() *auth.Client {
	return authClient
}

func ValidateToken(token, number string) error {
	if token == "" {
		return errors.NewError("Token not provided")
	}
	tokenCtx, err := AuthClient().VerifyIDToken(context.Background(), token)
	if err != nil {
		return err
	}
	phoneNumber := tokenCtx.Claims["phone_number"].(string)
	if phoneNumber[1:] != number {
		return errors.NewError("Phone number not matched with token")
	}
	return nil
}
