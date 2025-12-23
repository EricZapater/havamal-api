package versions

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(version *Version) error
	GetAll() ([]Version, error)
	GetById(id uuid.UUID) (*Version, error)
	Update(id uuid.UUID, version *Version) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(version *Version) error {
	query := `INSERT INTO versions (id, version, post_id, version_number, content, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, version.ID, version.Version, version.PostId, version.VersionNumber, version.Content, version.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAll() ([]Version, error) {
	query := `SELECT id, version, post_id, version_number, content, created_at FROM versions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var versions []Version
	for rows.Next() {
		var version Version
		if err := rows.Scan(&version.ID, &version.Version, &version.PostId, &version.VersionNumber, &version.Content, &version.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (r *repository) GetById(id uuid.UUID) (*Version, error) {
	query := `SELECT id, version, post_id, version_number, content, created_at FROM versions WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var version Version
	if err := row.Scan(&version.ID, &version.Version, &version.PostId, &version.VersionNumber, &version.Content, &version.CreatedAt); err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *repository) Update(id uuid.UUID, version *Version) error {
	query := `UPDATE versions SET version = $2, post_id = $3, version_number = $4, content = $5, created_at = $6 WHERE id = $1`
	_, err := r.db.Exec(query, id, version.Version, version.PostId, version.VersionNumber, version.Content, version.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(id uuid.UUID) error {
	query := `DELETE FROM versions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
