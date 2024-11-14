package auth

import "github.com/CelestialCrafter/lang-guesser/common/auth"

type Provider interface {
	GetUrl(state string) string
	Exchange(code string) (*auth.UserClaims, error)
}

var providers = map[string]Provider{}

func InitializeProviders() {
	providers = map[string]Provider{
		"google": NewGoogleProvider(),
	}
}

