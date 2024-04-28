package user

import (
	"database/sql"
	"log/slog"
	"strconv"
)

type userDb struct {
	ID    int32  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

type UserStorage struct {
	db            *sql.DB
	getUsersStm   *sql.Stmt
	createUserStm *sql.Stmt
	getUserStm    *sql.Stmt
}

func NewUserDb(db *sql.DB) *UserStorage {
	slog.Info("Initialising User Storage")

	getUsersStm, err := db.Prepare("SELECT * FROM users")
	if err != nil {
		return nil
	}

	createUserStm, err := db.Prepare("INSERT INTO users (name,email) VALUES ( ?, ?);")
	if err != nil {
		return nil
	}

	getUserStm, err := db.Prepare("SELECT * FROM users WHERE id = ?;")
	if err != nil {
		return nil
	}
	slog.Info("Intialised User Storage")

	return &UserStorage{db: db, getUsersStm: getUsersStm, getUserStm: getUserStm, createUserStm: createUserStm}
}

func (s *UserStorage) createUser(name, email string) (string, error) {
	res, err := s.createUserStm.Exec(name, email)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(id)), nil
}

func (s *UserStorage) getUsers() ([]userDb, error) {
	rows, err := s.getUsersStm.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]userDb, 3)

	for rows.Next() {
		var user userDb
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) getUser(id int) (*userDb, error) {
	var user userDb
	err := s.getUserStm.QueryRow(id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
