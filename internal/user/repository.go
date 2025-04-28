package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserRepository struct {
	db *pgxpool.Pool
}

func NewRepository(c *pgxpool.Pool) Repository {
	return &postgresUserRepository{db: c}
}

func (p *postgresUserRepository) CreateOne(user *User) error {
	sql := `INSERT INTO users (email, password_hash, is_superuser) VALUES ($1, $2, $3)`
	_, err := p.db.Exec(context.Background(), sql, user.Email, user.PasswordHash, user.IsSuperuser)
	return err
}

func (p *postgresUserRepository) UpdateOneById(id string, user *User) error {
	return nil
}

func (p *postgresUserRepository) DeleteOneById(id string) error {
	return nil
}

func (p *postgresUserRepository) GetOneByEmail(email string) (*User, error) {
	var user User
	err := p.db.QueryRow(context.Background(), `SELECT * FROM users WHERE email=$1`, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsSuperuser,
	)
	return &user, err
}

func (p *postgresUserRepository) GetOneById(id string) (*User, error) {
	return &User{}, nil
}

func (p *postgresUserRepository) GetAll() ([]User, error) {
	return []User{}, nil
}
