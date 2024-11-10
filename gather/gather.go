package gather

import (
	"context"
	"math"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/CelestialCrafter/lang-guesser/common"
	"github.com/CelestialCrafter/lang-guesser/ratelimit"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v66/github"
)

var (
	blobAmount = 50
	targetKb = 20
	minStars = 200
)

var suffixes = map[string]string{
	"go": "go",
	"rust": "rs",
	"python": "py",
}

type repository struct {
	Name string
	Owner string
}

type blob struct {
	SHA string
	Path string
	Size int
}

func SortBySize(blobs []blob, optimal int) {
	slices.SortFunc(blobs, func(a blob, b blob) int {
		return int(math.Abs(float64(optimal - a.Size)) - math.Abs(float64(optimal - b.Size)))
	})
}

func FilterBySuffix(blobs []blob, suffix string) []blob {
	return slices.DeleteFunc(blobs, func(b blob) bool {
		return !strings.HasSuffix(b.Path, suffix)
	})
}

func Gather() {
	_, ok := suffixes[*common.Gather]
	if !ok {
		log.Fatal("language not supported")
	}

	ratelimit.ConcurrentPermits.Aquire()
	defer ratelimit.ConcurrentPermits.Release()

	token, ok := os.LookupEnv("GITHUB_TOKEN") 
	if !ok {
		log.Fatal("GITHUB_TOKEN not set")
	}

	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()

	// select repo
	repos, err := GetRepos(ctx, client, *common.Gather, minStars)
	if err != nil {
		log.Fatal("could not get repositories", "error", err)
	}
	repo := repos[rand.IntN(len(repos))]
	log.Info("using repository", "repository", repo)

	// get blobs from default tree
	branch, err := GetDefaultBranch(ctx, client, repo)
	if err != nil {
		log.Fatal("could not get default branch", "error", err)
	}

	blobs, err := GetTree(ctx, client, repo, branch)
	if err != nil {
		log.Fatal("could not fetch blobs from repo", "error", err, "repository", repo)
	}

	// filter and sort
	blobs = FilterBySuffix(blobs, suffixes[*common.Gather])
	SortBySize(blobs, targetKb * 250)

	// download and parse blobs
	if len(blobs) > blobAmount {
		blobs = blobs[:blobAmount]
	}

	wg := sync.WaitGroup{}
	for i := range blobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer ratelimit.ConcurrentPermits.Release()
			ratelimit.ConcurrentPermits.Aquire()

			err := DownloadBlob(ctx, client, repo, blobs[i])
			if err != nil {
				log.Error("could not download blob", "error", err)
			}
		}()
	}

	wg.Wait()
}

