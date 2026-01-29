package navigation

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(navigation *Navigation) error
	GetAll() ([]Navigation, error)
	GetById(id uuid.UUID) (*Navigation, error)
	GetBySlug(slug string) (*Navigation, error)
	Update(id uuid.UUID, navigation *Navigation) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(navigation *Navigation) error {
	query := `INSERT INTO navigation (id, label, slug, type, "order", parent_id, link_source, category_id, post_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	var parentID, categoryID, postID interface{}
	if navigation.ParentId != nil {
		parentID = *navigation.ParentId
	}
	if navigation.CategoryId != nil {
		categoryID = *navigation.CategoryId
	}
	if navigation.PostId != nil {
		postID = *navigation.PostId
	}
	_, err := r.db.Exec(query, navigation.ID, navigation.Label, navigation.Slug, navigation.Type, navigation.Order, parentID, navigation.LinkSource, categoryID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAll() ([]Navigation, error) {
	query := `SELECT id, label, slug, type, "order", parent_id, link_source, category_id, post_id FROM navigation`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var navigations []Navigation
	for rows.Next() {
		var navigation Navigation
		var parentID, categoryID, postID uuid.NullUUID
		if err := rows.Scan(&navigation.ID, &navigation.Label, &navigation.Slug, &navigation.Type, &navigation.Order, &parentID, &navigation.LinkSource, &categoryID, &postID); err != nil {
			return nil, err
		}
		if parentID.Valid {
			navigation.ParentId = &parentID.UUID
		}
		if categoryID.Valid {
			navigation.CategoryId = &categoryID.UUID
		}
		if postID.Valid {
			navigation.PostId = &postID.UUID
		}
		navigations = append(navigations, navigation)
	}
	return navigations, nil
}

func (r *repository) GetById(id uuid.UUID) (*Navigation, error) {
	query := `SELECT id, label, slug, type, "order", parent_id, link_source, category_id, post_id FROM navigation WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var navigation Navigation
	var parentID, categoryID, postID uuid.NullUUID
	if err := row.Scan(&navigation.ID, &navigation.Label, &navigation.Slug, &navigation.Type, &navigation.Order, &parentID, &navigation.LinkSource, &categoryID, &postID); err != nil {
		return nil, err
	}
	if parentID.Valid {
		navigation.ParentId = &parentID.UUID
	}
	if categoryID.Valid {
		navigation.CategoryId = &categoryID.UUID
	}
	if postID.Valid {
		navigation.PostId = &postID.UUID
	}
	return &navigation, nil
}

func (r *repository) GetBySlug(slug string) (*Navigation, error) {
	query := `SELECT id, label, slug, type, "order", parent_id, link_source, category_id, post_id FROM navigation WHERE slug = $1`
	row := r.db.QueryRow(query, slug)
	var navigation Navigation
	var parentID, categoryID, postID uuid.NullUUID
	if err := row.Scan(&navigation.ID, &navigation.Label, &navigation.Slug, &navigation.Type, &navigation.Order, &parentID, &navigation.LinkSource, &categoryID, &postID); err != nil {
		return nil, err
	}
	if parentID.Valid {
		navigation.ParentId = &parentID.UUID
	}
	if categoryID.Valid {
		navigation.CategoryId = &categoryID.UUID
	}
	if postID.Valid {
		navigation.PostId = &postID.UUID
	}
	return &navigation, nil
}

func (r *repository) Update(id uuid.UUID, navigation *Navigation) error {
	query := `UPDATE navigation SET label = $2, slug = $3, type = $4, "order" = $5, parent_id = $6, link_source = $7, category_id = $8, post_id = $9 WHERE id = $1`
	var parentID, categoryID, postID interface{}
	if navigation.ParentId != nil {
		parentID = *navigation.ParentId
	}
	if navigation.CategoryId != nil {
		categoryID = *navigation.CategoryId
	}
	if navigation.PostId != nil {
		postID = *navigation.PostId
	}
	_, err := r.db.Exec(query, id, navigation.Label, navigation.Slug, navigation.Type, navigation.Order, parentID, navigation.LinkSource, categoryID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(id uuid.UUID) error {
	query := `DELETE FROM navigation WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
