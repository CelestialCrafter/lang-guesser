package gather

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v66/github"
)

type repository struct {
	Name string
	Owner string
}

func GetRepos(ctx context.Context, client *github.Client, language string, minStars int) ([]repository, error) {
	opts := &github.SearchOptions {
		Sort: "stars",
	}

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

	return repos, nil
}

type blob struct {
	SHA string
	Path string
	Size int
}

func GetBlobs(ctx context.Context, client *github.Client, repo repository) ([]blob, error) {
	data, _, err := client.Git.GetTree(ctx, repo.Owner, repo.Name, "master", true)
	if err != nil {
		return nil, err
	}

	blobs := make([]blob, len(data.Entries))
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

	return blobs, nil
}

func SortBySize(blobs []blob, optimal int) {
	slices.SortFunc(blobs, func(a blob, b blob) int {
		return int(math.Abs(float64(optimal - a.Size)) - math.Abs(float64(optimal - b.Size)))
	})
}

func FilterBySuffix(blobs []blob, suffix string) {
	slices.DeleteFunc(blobs, func(b blob) bool {
		return !strings.HasSuffix(b.Path, suffix)
	})
}

var (
	k = 100
	targetKb = 20
	minStars = 1000
	language = "rust"
	languageSuffix = "rs"
)

func Gather() {
	client := github.NewClient(nil)
	repos, err := GetRepos(context.Background(), client, language, minStars)
	if err != nil {
		log.Fatal("could not get repositories", "error", err)
	}

	repo := repos[0]
	blobs, err := GetBlobs(context.Background(), client, repo)
	if err != nil {
		log.Fatal("could not fetch blobs from repo", "error", err, "repository", repo)
	}

	FilterBySuffix(blobs, languageSuffix)
	SortBySize(blobs, targetKb * 250)
	blob := blobs[rand.IntN(k)]

	data, _, err := client.Git.GetBlobRaw(context.Background(), repo.Owner, repo.Name, blob.SHA)
	if err != nil {
		log.Fatal("could not fetch blob data", "error", err)
	}

	log.Info("downloading file", "file", blob.Path)
	file, err := os.Create(path.Base(blob.Path))
	if err != nil {
		log.Fatal("could not create file", "error", err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatal("could not write to file", "error", err)
	}
}
