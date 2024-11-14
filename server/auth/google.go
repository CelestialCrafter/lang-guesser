package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/CelestialCrafter/lang-guesser/common/auth"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const redirectUrl = "http://localhost:8080/api/auth/google/callback"

type Google struct {
	config *oauth2.Config
}

type googleClaims struct {
	ID string `json:"id"`
	Picture string `json:"picture"`
	Name string `json:"name"`
}

func NewGoogleProvider() Google {
	return Google {
		config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{
				"openid",
			 	"https://www.googleapis.com/auth/userinfo.email",
			 	"https://www.googleapis.com/auth/userinfo.profile",
			},
			RedirectURL:  redirectUrl,
			Endpoint:     google.Endpoint,
		},
	}
}


func (g Google) GetUrl(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g Google) Exchange(code string) (*auth.UserClaims, error) {
	oauthToken, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange: %w", err)
	}

	client := g.config.Client(context.Background(), oauthToken)
	res, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return nil, fmt.Errorf("fetch user: %w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// sign token
	googleClaims := new(googleClaims)
	err = json.Unmarshal(body, googleClaims)
	if err != nil {
		return nil, err
	}

	claims := &auth.UserClaims{
		ID: fmt.Sprint("google-", googleClaims.ID),
		Username: googleClaims.Name,
		Picture: googleClaims.Picture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	return claims, nil
}

