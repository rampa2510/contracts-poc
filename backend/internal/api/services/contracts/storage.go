package contracts

import (
	"database/sql"
	"log/slog"
	"strconv"

	"github.com/rampa2510/contracts-poc/config"
)

type contractsDb struct {
	S3Key  string `json:"s3Key" db:"s3_key"`
	UserId int    `json:"userId" db:"user_id"`
}

type ContractsStorage struct {
	db                 *sql.DB
	createContractsStm *sql.Stmt
	env                *config.EnvVars
}

func NewContractsStorage(db *sql.DB) *ContractsStorage {
	slog.Info("Intialising Contracts Storage")

	createContractsStm, err := db.Prepare("INSERT INTO contracts (s3_key,user_id) VALUES (?, ?);")
	if err != nil {
		return nil
	}

	slog.Info("Initialised Contracts Storage")
	return &ContractsStorage{db: db, createContractsStm: createContractsStm}
}

func (cs *ContractsStorage) CreateNewContract(s3Key string, userId int) (string, error) {
	res, err := cs.createContractsStm.Exec(s3Key, userId)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(id)), nil
}
