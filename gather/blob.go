package gather

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/CelestialCrafter/lang-guesser/common"
	"github.com/CelestialCrafter/lang-guesser/db"
	"github.com/CelestialCrafter/lang-guesser/ratelimit"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v66/github"
)

func ASTSplitProtocol(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	index := bytes.Index(data, []byte("|"))

	// gather more bytes for header
	if index == -1 {
		return 0, nil, nil
	}

	length, err := strconv.Atoi(string(data[:index]))
	if err != nil {
		return 0, nil, err
	}
	index++

	if len(data) < index + length {
		return 0, nil, nil
	}

	return index + length, data[index:index + length], nil
}

func ParseBlob(data []byte) ([][]byte, error) {
	cmd := exec.Command("sh", "start.sh")

	cmd.Dir = path.Join("ast-processors", *common.Gather) 
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

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}


	return sections, nil
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

	if *common.Testcases && len(sections) > 0 {
		basePath := path.Join("testcases", blob.SHA)
		err := os.Mkdir(basePath, os.ModePerm)
		if err != nil {
			return err
		}

		mainPath := fmt.Sprintf("%s.%s", "main", LangToSuffix[*common.Gather]) 
		err = os.WriteFile(path.Join(basePath, mainPath), data, 0o666)
		if err != nil {
			return err
		}

		for i, section := range sections {
			sectionPath := fmt.Sprintf("%s.%s", strconv.Itoa(i), LangToSuffix[*common.Gather])
			err := os.WriteFile(path.Join(basePath, sectionPath), section, 0o666)
			if err != nil {
				log.Error("could not save file", "error", err)
			}
		}

		return nil
	}

	for _, section := range sections {
		err = db.CreateChallenge(db.Challenge{
			Sha: blob.SHA,
			Code: section,
			Language: *common.Gather,
		})
		if err != nil {
			return err
		}

		log.Info("downloaded blob", "sha", blob.SHA)
	}

	return nil
}

