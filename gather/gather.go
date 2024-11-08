package gather

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path"
	"slices"
	"strings"
	"sync"

	"github.com/CelestialCrafter/lang-guesser/ratelimit"
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

type blob struct {
	SHA string
	Path string
	Size int
}

func GetBlobs(ctx context.Context, client *github.Client, repo repository, branch string) ([]blob, error) {
	ratelimit.EndpointPermits[ratelimit.GetTree].Aquire()
	data, _, err := client.Git.GetTree(ctx, repo.Owner, repo.Name, branch, true)
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

func DownloadBlob(ctx context.Context, client *github.Client, repo repository, blob blob, dir string) error {
	ratelimit.EndpointPermits[ratelimit.GetBlob].Aquire()
	data, _, err := client.Git.GetBlobRaw(context.Background(), repo.Owner, repo.Name, blob.SHA)
	if err != nil {
		return err
	}

	filepath := path.Join(dir, path.Base(blob.Path)) 
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	log.Info("downloaded blob", "path", filepath)
	
	return nil
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
	blobAmount = 50
	targetKb = 20
	minStars = 1000
	language = "rust"
	languageSuffix = "rs"
)

func Gather() {
	ratelimit.ConcurrentPermits.Aquire()
	defer ratelimit.ConcurrentPermits.Release()

	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	ctx := context.Background()

	repos, err := GetRepos(ctx, client, language, minStars)
	if err != nil {
		log.Fatal("could not get repositories", "error", err)
	}

	repo := repos[rand.IntN(len(repos))]
	log.Info("using repository", "repository", repo)

	branch, err := GetDefaultBranch(ctx, client, repo)
	if err != nil {
		log.Fatal("could not get default branch", "error", err)
	}

	blobs, err := GetBlobs(ctx, client, repo, branch)
	if err != nil {
		log.Fatal("could not fetch blobs from repo", "error", err, "repository", repo)
	}

	FilterBySuffix(blobs, languageSuffix)
	SortBySize(blobs, targetKb * 250)

	wg := sync.WaitGroup{}
	for i := range blobs[:blobAmount] {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer ratelimit.ConcurrentPermits.Release()
			ratelimit.ConcurrentPermits.Aquire()

			DownloadBlob(ctx, client, repo, blobs[i], "files")
		}()
	}

	wg.Wait()
}

