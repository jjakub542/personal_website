package contact

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresMessageFormRepository struct {
	db *pgxpool.Pool
}

func NewRepository(c *pgxpool.Pool) Repository {
	return &postgresMessageFormRepository{db: c}
}

func (p *postgresMessageFormRepository) CreateOne(c *MessageForm) error {
	sql := `INSERT INTO messages (name, email, message) VALUES ($1, $2, $3)`
	_, err := p.db.Exec(context.Background(), sql, c.Name, c.Email, c.Message)
	return err
}

func (p *postgresMessageFormRepository) DeleteOneById(id string) error {
	sql := `DELETE FROM messages WHERE id=$1`
	_, err := p.db.Exec(context.Background(), sql, id)
	return err
}

func (p *postgresMessageFormRepository) GetAll() ([]MessageForm, error) {
	var msgs []MessageForm
	sql := `SELECT m.id, m.name, m.email, m.message,
        m.created_at FROM messages m ORDER BY created_at DESC;`
	rows, err := p.db.Query(context.Background(), sql)
	if err != nil {
		return msgs, err
	}
	for rows.Next() {
		var msg MessageForm
		if err := rows.Scan(
			&msg.Id,
			&msg.Name,
			&msg.Email,
			&msg.Message,
			&msg.CreatedAt,
		); err != nil {
			return msgs, err
		}
		msgs = append(msgs, msg)
	}
	if err = rows.Err(); err != nil {
		return msgs, err
	}
	return msgs, nil
}
