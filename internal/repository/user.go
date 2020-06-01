package repository

import (
	"adeia-api/internal/model"
	"adeia-api/internal/service/db"

	"github.com/jmoiron/sqlx"
)

const (
	queryUserInsert  = "INSERT INTO users (employee_id, name, email) VALUES (:employee_id, :name, :email)"
	queryUserByEmpID = "SELECT * FROM users WHERE employee_id=$1"
)

// PGUserRepo is an implementation of UserRepo for Postgres.
type PGUserRepo struct {
	conn *db.DB
}

// NewPGUserRepo creates a new PGUserRepo.
func NewPGUserRepo() UserRepo {
	c := db.GetConn()
	return &PGUserRepo{c}
}

// InsertWithTx inserts a user using the provided transaction.
func (p *PGUserRepo) InsertWithTx(tx *sqlx.Tx, u *model.User) error {
	_, err := tx.NamedExec(queryUserInsert, u)
	return err
}

// Insert inserts a user using the db connection instance.
func (p *PGUserRepo) Insert(u *model.User) error {
	_, err := p.conn.NamedExec(queryUserInsert, u)
	return err
}

// GetByEmpID gets a user from db using the empId.
func (p *PGUserRepo) GetByEmpID(empID string) (*model.User, error) {
	u := model.User{}
	err := p.conn.Get(&u, queryUserByEmpID, empID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
