package gather

import (
	"context"

	"github.com/CelestialCrafter/lang-guesser/db"
	"github.com/CelestialCrafter/lang-guesser/ratelimit"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v66/github"
)

func DownloadBlob(ctx context.Context, client *github.Client, language string, repo repository, blob blob) error {
	ratelimit.EndpointPermits[ratelimit.GetBlob].Aquire()
	data, _, err := client.Git.GetBlobRaw(context.Background(), repo.Owner, repo.Name, blob.SHA)
	if err != nil {
		return err
	}

	sections, err := ParseSections(data, language)
	if err != nil {
		return err
	}

	for _, section := range sections {
		err = db.CreateChallenge(db.Challenge{
			Sha: blob.SHA,
			Code: section,
			Language: language,
		})
		if err != nil {
			return err
		}

		log.Info("downloaded blob", "sha", blob.SHA)
	}

	return nil
}

