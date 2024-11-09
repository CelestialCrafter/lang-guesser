package gather

import (
	"context"
	"fmt"

	"github.com/CelestialCrafter/lang-guesser/ratelimit"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v66/github"
)

func GetRepos(ctx context.Context, client *github.Client, language string, minStars int) ([]repository, error) {
	opts := &github.SearchOptions {
		Sort: "updated",
	}

	ratelimit.EndpointPermits[ratelimit.Search].Aquire()
	data, _, err := client.Search.Repositories(ctx, fmt.Sprintf("stars:>=%d language:%s", minStars, language), opts)
	if err != nil {
		return nil, err
	}

	// @TODO filter by license?
	repos := make([]repository, len(data.Repositories))
	for i, entry := range data.Repositories {
		if entry.Name == nil || entry.Owner == nil || entry.Owner.Login == nil {
			continue
		}

		repos[i] = repository{
			Name: *entry.Name,
			Owner: *entry.Owner.Login,
		}
	}

	log.Info("fetched repos", "amount", len(repos))

	return repos, nil
}

