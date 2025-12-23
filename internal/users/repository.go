package users

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, id uuid.UUID) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user User) (User, error) {
	query := `INSERT INTO users (id, username, email, password, 
						is_admin, is_active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password, 
		user.IsAdmin, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) FindAll(ctx context.Context) ([]User, error) {
	query := `SELECT id, username, email, password, is_admin, is_active, created_at, updated_at
				FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (User, error) {
	query := `SELECT id, username, email, password, is_admin, is_active, created_at, updated_at 
				FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (User, error) {
	query := `SELECT id, username, email, password, is_admin, is_active, created_at, updated_at
				FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) Update(ctx context.Context, user User) (User, error) {
	query := `UPDATE users SET username = $2, email = $3, password = $4, is_admin = $5, is_active = $6, updated_at = $7 WHERE id = $1 RETURNING id`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password, user.IsAdmin, user.IsActive, user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
