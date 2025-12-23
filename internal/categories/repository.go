package categories

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(category *Category) error
	GetAll() ([]Category, error)
	GetById(id uuid.UUID) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	Update(category *Category) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(category *Category) error {
	query := `INSERT INTO categories (id, name, slug, description, "order", created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, category.ID, category.Name, category.Slug, category.Description, category.Order, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAll() ([]Category, error) {
	query := `SELECT id, name, slug, description, "order", created_at, updated_at
	FROM categories`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug, &category.Description, &category.Order, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *repository) GetById(id uuid.UUID) (*Category, error) {
	query := `SELECT id, name, slug, description, "order", created_at, updated_at
	FROM categories
	WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var category Category
	if err := row.Scan(&category.ID, &category.Name, &category.Slug, &category.Description, &category.Order, &category.CreatedAt, &category.UpdatedAt); err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) GetBySlug(slug string) (*Category, error) {
	query := `SELECT id, name, slug, description, "order", created_at, updated_at
	FROM categories
	WHERE slug = $1`
	row := r.db.QueryRow(query, slug)
	var category Category
	if err := row.Scan(&category.ID, &category.Name, &category.Slug, &category.Description, &category.Order, &category.CreatedAt, &category.UpdatedAt); err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) Update(category *Category) error {
	query := `UPDATE categories
	SET name = $2, slug = $3, description = $4, "order" = $5, updated_at = $6
	WHERE id = $1`
	_, err := r.db.Exec(query, category.ID, category.Name, category.Slug, category.Description, category.Order, category.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
