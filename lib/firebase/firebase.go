package firebase

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/techartificer/swiftex/lib/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var authClient *auth.Client

func Initialize() error {
	// log.Println(os.Getenv("FIREBASE"))
	creds, err := google.CredentialsFromJSON(context.Background(), []byte(os.Getenv("FIREBASE")))
	if err != nil {
		return err
	}
	opt := option.WithCredentials(creds)
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
