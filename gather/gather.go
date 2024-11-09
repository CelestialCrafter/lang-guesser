package gather

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"os/exec"
	"path"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/CelestialCrafter/lang-guesser/db"
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

func ASTSplitProtocol(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	lengthBytes, section, found := bytes.Cut(data, []byte("|"))

	// gather more bytes for header
	if !found {
		return 0, nil, nil
	}

	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return 0, nil, err
	}

	// gather more data
	if !found || len(section) != length {
		return 0, nil, nil
	}

	return len(lengthBytes) + length, section, nil
}

func ParseBlob(data []byte) ([][]byte, error) {
	cmd := exec.Command("sh", "start.sh")

	cmd.Dir = path.Join("ast-processors", language) 
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	_, err = stdin.Write(data)
	if err != nil {
		return nil, err
	}

	err = stdin.Close()
	if err != nil {
		return nil, err
	}

	// len(data) is used as an initial size for each section, with 5 sections by default)
	sections := make([][]byte, 0)

	scanner := bufio.NewScanner(stdout)
	scanner.Split(ASTSplitProtocol)

	for scanner.Scan() {
		if scanner.Err() != nil {
			return nil, err
		}

		sections = append(sections, scanner.Bytes())
	}

	return sections, cmd.Wait()
}

func DownloadBlob(ctx context.Context, client *github.Client, repo repository, blob blob) error {
	ratelimit.EndpointPermits[ratelimit.GetBlob].Aquire()
	data, _, err := client.Git.GetBlobRaw(context.Background(), repo.Owner, repo.Name, blob.SHA)
	if err != nil {
		return err
	}

	sections, err := ParseBlob(data)
	if err != nil {
		return err
	}

	for _, data := range sections {
		err = db.CreateChallenge(db.Challenge{
			Sha: blob.SHA,
			Code: data,
			Language: language,
		})
		if err != nil {
			return err
		}

		log.Info("downloaded blob", "sha", blob.SHA)
	}

	return nil
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

	token, ok := os.LookupEnv("GITHUB_TOKEN") 
	if !ok {
		log.Fatal("GITHUB_TOKEN not set")
	}

	client := github.NewClient(nil).WithAuthToken(token)
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

	blobs = FilterBySuffix(blobs, languageSuffix)
	SortBySize(blobs, targetKb * 250)


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

