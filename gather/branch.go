package gather

import (
	"context"
	"errors"
	"slices"

	"github.com/CelestialCrafter/lang-guesser/ratelimit"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v66/github"
)

func GetTree(ctx context.Context, client *github.Client, repo repository, branch string) ([]blob, error) {
	ratelimit.EndpointPermits[ratelimit.GetTree].Aquire()
	data, _, err := client.Git.GetTree(ctx, repo.Owner, repo.Name, branch, true)
	if err != nil {
		return nil, err
	}

	blobs := make([]blob, 0, len(data.Entries))
	for _, entry := range data.Entries {
		if entry.GetType() != "blob" || entry.SHA == nil || entry.Path == nil {
			continue
		}

		blobs = append(blobs, blob{
			SHA: *entry.SHA,
			Path: *entry.Path,
			Size: *entry.Size,
		})
	}
	blobs = slices.Clip(blobs)

	log.Info("fetched blobs", "amount", len(blobs))

	return blobs, nil
}

func GetDefaultBranch(ctx context.Context, client *github.Client, repo repository) (string, error) {
	ratelimit.EndpointPermits[ratelimit.GetRepo].Aquire()
	data, _, err := client.Repositories.Get(ctx, repo.Owner, repo.Name)
	if err != nil {
		return "", err
	}

	if data.DefaultBranch == nil {
		return "", errors.New("no default branch")
	}

	return *data.DefaultBranch, nil
}

