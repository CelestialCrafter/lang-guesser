package gather_test

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/CelestialCrafter/lang-guesser/common"
	"github.com/CelestialCrafter/lang-guesser/gather"
)

func TestTestcases(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	base := "testcases"
	testcases, err := os.ReadDir(base)
	if err != nil {
		panic(err)
	}

	for _, testcase := range testcases {
		testTestcase(t, path.Join(base, testcase.Name()))
	}
}

func testTestcase(t *testing.T, base string) {
	files, err := os.ReadDir(base)
	if err != nil {
		panic(err)
	}

	var suffix string
	var main []byte
	expectedSections := make(map[[20]byte]struct{}, len(files) - 1)

	for _, file := range files {
		path := path.Join(base, file.Name())
		if suffix == "" {
			suffix = strings.Split(file.Name(), ".")[1]
			if gather.LangToSuffix[*common.Gather] != suffix {
				return
			}
		}

		data, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		if file.Name() == "main." + suffix {
			main = data
		} else {
			hash := sha1.Sum(data)
			expectedSections[hash] = struct{}{}
		}
	}

	sections, err := gather.ParseBlob(main)
	if err != nil {
		panic(err)
	}

	for _, section := range sections {
		hash := sha1.Sum(section)
		hex := hex.EncodeToString(hash[:])
		_, ok := expectedSections[hash]
		if !ok {
			t.Errorf("section failed checksum. got: %s", hex)
		}
	}
}
