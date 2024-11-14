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
)

type Github struct {
	config *oauth2.Config
}

type githubClaims struct {
	ID int `json:"id"`
	Avatar string `json:"avatar_url"`
	Name string `json:"name"`
}

func NewGithubProvider() Github {
	return Github {
		config: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			Scopes:       []string{},
			RedirectURL:  "http://localhost:8080/api/auth/github/callback",
			Endpoint:     oauth2.Endpoint{
				AuthURL: "https://github.com/login/oauth/authorize",
				DeviceAuthURL: "https://github.com/login/device/code",
				TokenURL: "https://github.com/login/oauth/access_token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
	}
}


func (g Github) GetUrl(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g Github) Exchange(code string) (*auth.UserClaims, error) {
	oauthToken, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange: %w", err)
	}

	client := g.config.Client(context.Background(), oauthToken)
	res, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("fetch user: %w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// sign token
	githubClaims := new(githubClaims)
	err = json.Unmarshal(body, githubClaims)
	if err != nil {
		return nil, err
	}

	claims := &auth.UserClaims{
		ID: fmt.Sprint("google-", githubClaims.ID),
		Username: githubClaims.Name,
		Picture: githubClaims.Avatar,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	return claims, nil
}

