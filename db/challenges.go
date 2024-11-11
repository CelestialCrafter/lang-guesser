package db

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const challengeDbPath = "challenges.db"

var challenges *sqlx.DB

const MIN_LENGTH = 250

func InitChallenges() {
	if _, err := os.Stat(challengeDbPath); os.IsNotExist(err) {
		file, err := os.Create(challengeDbPath)
		if err != nil {
			log.Fatal("could not create database file", "error", err)
		}
		file.Close()
	}

	var err error
	challenges, err = sqlx.Connect("sqlite3", challengeDbPath)
	if err != nil {
		log.Fatal("could not open database", "error", err)
	}

	_, err = challenges.Exec(`CREATE TABLE IF NOT EXISTS challenges (
		sha TEXT NOT NULL,
		code BLOB PRIMARY KEY,
		language TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal("could not create challenges table", "error", err)
	}

	log.Info("initialized database")
}

type Challenge struct {
	Sha string `json:"sha"`
	Code []byte `json:"code"`
	Language string `json:"language"`
}

func CreateChallenge(challenge Challenge) error {
	_, err := challenges.Exec("INSERT INTO challenges (sha, code, language) VALUES (?, ?, ?)", challenge.Sha, challenge.Code, challenge.Language)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

func GetRandomChallenge() (*Challenge, error) {
	challenge := new(Challenge)
	err := challenges.Get(challenge, fmt.Sprintf("SELECT * FROM challenges WHERE length(code) > %d ORDER BY RANDOM() LIMIT 1;", MIN_LENGTH))
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return challenge, nil
}
