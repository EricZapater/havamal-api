package posts

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	CreatePost(post *Post) error
	GetPost(id uuid.UUID) (*Response, error)	
	GetPosts() ([]Response, error)
	GetPublishedPosts() ([]Response, error) 
	GetPostsByAuthor(authorId uuid.UUID) ([]Response, error)
	GetPostBySlug(slug string) (*Response, error)
	GetSummariesByCategory(category string) ([]Response, error)
	UpdatePost(post *Post) error
	DeletePost(id uuid.UUID) error
	AddCategory(request PostCategories) error
	DeleteCategory(request PostCategories) error
	AddVersion(request PostVersion) error
	DeleteVersion(request PostVersion) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreatePost(post *Post) error {
	query := `INSERT INTO posts (id, title, slug, summary, content, status, published_at, updated_at, author_id, columns)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`	
	_, err := r.db.Exec(query, post.ID, post.Title, post.Slug, post.Summary, post.Content, post.Status, post.PublishedAt, post.UpdatedAt, post.AuthorId, post.Columns)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetPost(id uuid.UUID) (*Response, error) {
	query := `SELECT p.id, p.title, p.slug, p.summary, p.content, p.status, p.published_at, 
					p.updated_at, p.author_id, pc.category_id, c.name as category_name, c.description as category_description, 
					c.slug as category_slug, u.username as author_name, p.columns
	FROM posts p
		INNER JOIN post_categories pc ON p.Id = pc.post_id
		INNER JOIN categories c ON pc.category_id = c.id
		INNER JOIN users u ON p.author_id = u.id
	WHERE p.id = $1`
	row := r.db.QueryRow(query, id)
	var post Response
	if err := row.Scan(&post.ID, &post.Title, &post.Slug, &post.Summary, &post.Content, &post.Status, &post.PublishedAt, 
						&post.UpdatedAt, &post.AuthorId, &post.CategoryId, &post.CategoryName, &post.CategoryDescription, 
						&post.CategorySlug, &post.AuthorName, &post.Columns); err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) GetPosts() ([]Response, error) {
	query := `SELECT p.id, p.title, p.slug, p.summary, p.content, p.status, p.published_at, 
					p.updated_at, p.author_id, pc.category_id, c.name as category_name, c.description as category_description, 
					c.slug as category_slug, u.username as author_name, p.columns
	FROM posts p
		INNER JOIN post_categories pc ON p.Id = pc.post_id
		INNER JOIN categories c ON pc.category_id = c.id
		INNER JOIN users u ON p.author_id = u.id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Response
	for rows.Next() {
		var post Response
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Summary, &post.Content, &post.Status, &post.PublishedAt, 
							&post.UpdatedAt, &post.AuthorId, &post.CategoryId, &post.CategoryName, &post.CategoryDescription, 
							&post.CategorySlug, &post.AuthorName, &post.Columns); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *repository) GetPublishedPosts() ([]Response, error) {
	query := `SELECT p.id, p.title, p.slug, p.summary, p.content, p.status, p.published_at, 
					p.updated_at, p.author_id, pc.category_id, c.name as category_name, c.description as category_description, 
					c.slug as category_slug, u.username as author_name, p.columns
	FROM posts p
		INNER JOIN post_categories pc ON p.Id = pc.post_id
		INNER JOIN categories c ON pc.category_id = c.id
		INNER JOIN users u ON p.author_id = u.id
	WHERE p.status = 'published'`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Response
	for rows.Next() {
		var post Response
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Summary, &post.Content, &post.Status, &post.PublishedAt, 
							&post.UpdatedAt, &post.AuthorId, &post.CategoryId, &post.CategoryName, &post.CategoryDescription, 
							&post.CategorySlug, &post.AuthorName, &post.Columns); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *repository) GetPostsByAuthor(authorId uuid.UUID) ([]Response, error) {
	query := `SELECT p.id, p.title, p.slug, p.summary, p.content, p.status, p.published_at, 
					p.updated_at, p.author_id, pc.category_id, c.name as category_name, c.description as category_description, c.slug as category_slug,
					u.username, p.columns
	FROM posts p
		INNER JOIN post_categories pc ON p.Id = pc.post_id
		INNER JOIN categories c ON pc.category_id = c.id
		INNER JOIN users u ON p.author_id = u.id
	WHERE p.author_id = $1`
	rows, err := r.db.Query(query, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Response
	for rows.Next() {
		var post Response
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Summary, &post.Content, &post.Status, &post.PublishedAt, 
							&post.UpdatedAt, &post.AuthorId, &post.CategoryId, &post.CategoryName, &post.CategoryDescription, &post.CategorySlug,
							&post.AuthorName, &post.Columns); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *repository) GetPostBySlug(slug string) (*Response, error) {
	query := `SELECT p.id, p.title, p.slug, p.summary, p.content, p.status, p.published_at, 
					p.updated_at, p.author_id, pc.category_id, c.name as category_name, c.description as category_description, c.slug as category_slug,
					u.username, p.columns
	FROM posts p
		INNER JOIN post_categories pc ON p.Id = pc.post_id
		INNER JOIN categories c ON pc.category_id = c.id
		INNER JOIN users u ON p.author_id = u.id
	WHERE p.slug = $1`
	row := r.db.QueryRow(query, slug)
	var post Response
	if err := row.Scan(&post.ID, &post.Title, &post.Slug, &post.Summary, &post.Content, &post.Status, &post.PublishedAt, 
						&post.UpdatedAt, &post.AuthorId, &post.CategoryId, &post.CategoryName, &post.CategoryDescription, &post.CategorySlug,
						&post.AuthorName, &post.Columns); err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) GetSummariesByCategory(category string) ([]Response, error) {
	// Assuming category argument is likely a UUID string or slug. 
    // Try to match id first, but since the arg is string, standard implementation usually expects ID or uses separate queries.
    // Given the previous broken query 'WHERE category = $1', assuming ID matching on join for now.
	query := `SELECT p.id, p.title, p.slug, p.summary, p.content, p.status, p.published_at, 
					p.updated_at, p.author_id, pc.category_id, c.name as category_name, c.description as category_description, c.slug as category_slug,
					u.username, p.columns
	FROM posts p
		INNER JOIN post_categories pc ON p.Id = pc.post_id
		INNER JOIN categories c ON pc.category_id = c.id
		INNER JOIN users u ON p.author_id = u.id
	WHERE c.id::text = $1 OR c.slug = $1` 
	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Response
	for rows.Next() {
		var post Response
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Summary, &post.Content, &post.Status, &post.PublishedAt, 
							&post.UpdatedAt, &post.AuthorId, &post.CategoryId, &post.CategoryName, &post.CategoryDescription, &post.CategorySlug,
							&post.AuthorName, &post.Columns); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *repository) UpdatePost(post *Post) error {
	query := `UPDATE posts
	SET title = $2, slug = $3, summary = $4, content = $5, status = $6, published_at = $7, updated_at = $8, author_id = $9, columns = $10
	WHERE id = $1`
	_, err := r.db.Exec(query, post.ID, post.Title, post.Slug, post.Summary, post.Content, post.Status, post.PublishedAt, post.UpdatedAt, post.AuthorId, post.Columns)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeletePost(id uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) AddCategory(request PostCategories) error {
	query := `INSERT INTO post_categories (post_id, category_id)
	VALUES ($1, $2)`
	_, err := r.db.Exec(query, request.PostId, request.CategoryId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteCategory(request PostCategories) error {
	query := `DELETE FROM post_categories WHERE post_id = $1 AND category_id = $2`
	_, err := r.db.Exec(query, request.PostId, request.CategoryId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) AddVersion(request PostVersion) error {
	query := `INSERT INTO post_versions (version_id, post_id)
	VALUES ($1, $2)`
	_, err := r.db.Exec(query, request.VersionId, request.PostId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteVersion(request PostVersion) error {
	query := `DELETE FROM post_versions WHERE version_id = $1 AND post_id = $2`
	_, err := r.db.Exec(query, request.VersionId, request.PostId)
	if err != nil {
		return err
	}
	return nil
}
