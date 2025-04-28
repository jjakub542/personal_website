package article

import (
	"context"
	psql "database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresArticleRepository struct {
	db *pgxpool.Pool
}

func NewRepository(c *pgxpool.Pool) Repository {
	return &postgresArticleRepository{db: c}
}

func (p *postgresArticleRepository) GetAll() ([]Article, error) {
	var articles []Article
	sql := `SELECT a.id, a.title, a.description, a.content,
        a.created_at, a.updated_at, a.public FROM articles a ORDER BY created_at DESC;`
	rows, err := p.db.Query(context.Background(), sql)
	if err != nil {
		return articles, err
	}
	for rows.Next() {
		var article Article
		if err := rows.Scan(
			&article.Id,
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Public,
		); err != nil {
			return articles, err
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		return articles, err
	}
	return articles, nil
}

func (p *postgresArticleRepository) GetCount() (int, error) {
	var count int
	err := p.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM articles").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *postgresArticleRepository) GetAllPublicBetween(limit int, offset int) ([]Article, error) {
	var articles []Article
	sql := `SELECT a.id, a.title, a.description, a.content,
        a.created_at, a.updated_at, a.public, a.cover_image_id FROM articles a LEFT JOIN images i ON a.cover_image_id = i.id WHERE public=true ORDER BY created_at DESC LIMIT $1 OFFSET $2;`
	rows, err := p.db.Query(context.Background(), sql, limit, offset)
	if err != nil {
		return articles, err
	}
	defer rows.Close()
	for rows.Next() {
		var coverImgId psql.NullString
		var article Article
		if err := rows.Scan(
			&article.Id,
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Public,
			&coverImgId,
		); err != nil {
			return articles, err
		}
		if coverImgId.Valid {
			article.CoverImageId = &coverImgId.String
		} else {
			article.CoverImageId = nil
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		return articles, err
	}
	return articles, nil
}

func (p *postgresArticleRepository) GetAllPublic() ([]Article, error) {
	var articles []Article
	sql := `SELECT a.id, a.title, a.description, a.content,
        a.created_at, a.updated_at, a.public FROM articles a WHERE public=true ORDER BY created_at DESC;`
	rows, err := p.db.Query(context.Background(), sql)
	if err != nil {
		return articles, err
	}
	for rows.Next() {
		var article Article
		if err := rows.Scan(
			&article.Id,
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.Public,
		); err != nil {
			return articles, err
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		return articles, err
	}
	return articles, nil
}

func (p *postgresArticleRepository) GetOneById(id string) (*Article, error) {
	var article Article
	sql := `SELECT a.id, a.title, a.description, a.content,
        a.created_at, a.updated_at, a.public FROM articles a WHERE id=$1;`
	err := p.db.QueryRow(context.Background(), sql, id).Scan(
		&article.Id,
		&article.Title,
		&article.Description,
		&article.Content,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.Public,
	)
	return &article, err
}

func (p *postgresArticleRepository) CreateOne(a *Article) error {
	sql := `INSERT INTO articles (title, description, content, public)
	VALUES ($1, $2, $3, $4)`
	_, err := p.db.Exec(context.Background(), sql, a.Title, a.Description, a.Content, a.Public)
	return err
}

func (p *postgresArticleRepository) UpdateOneById(a *Article, id string) error {
	sql := `UPDATE articles SET
	title=$1, description=$2, content=$3, updated_at=$4, public=$5
	WHERE id=$6`
	_, err := p.db.Exec(context.Background(), sql, a.Title, a.Description, a.Content, time.Now(), a.Public, id)
	return err
}

func (p *postgresArticleRepository) DeleteOneById(id string) error {
	sql := `DELETE FROM articles WHERE id=$1`
	_, err := p.db.Exec(context.Background(), sql, id)
	return err
}

func (p *postgresArticleRepository) AttachImage(article_id string) (string, error) {
	var newImageID string
	sql := `INSERT INTO images (article_id) VALUES ($1) RETURNING id`
	err := p.db.QueryRow(context.Background(), sql, article_id).Scan(&newImageID)
	return newImageID, err
}

func (p *postgresArticleRepository) RemoveImage(filename string) error {
	// filename should be indexed !!!
	sql := `DELETE FROM images WHERE id=$1`
	_, err := p.db.Exec(context.Background(), sql, filename)
	return err
}

func (p *postgresArticleRepository) GetArticleImages(article_id string) ([]Image, error) {
	var images []Image
	sql := `SELECT * FROM images WHERE article_id=$1`
	rows, err := p.db.Query(context.Background(), sql, article_id)
	if err != nil {
		return images, err
	}
	for rows.Next() {
		var image Image
		if err := rows.Scan(
			&image.Id,
			&image.UploadedAt,
			&image.ArticleId,
		); err != nil {
			return images, err
		}
		images = append(images, image)
	}
	if err = rows.Err(); err != nil {
		return images, err
	}
	return images, nil
}

func (p *postgresArticleRepository) UpdateArticleCoverImage(article_id string, image_id string) error {
	fmt.Println("\n" + article_id + "\n")
	fmt.Println("\n" + image_id + "\n")
	sql := `UPDATE articles SET cover_image_id = $1 WHERE id = $2`
	_, err := p.db.Exec(context.Background(), sql, image_id, article_id)
	return err
}
