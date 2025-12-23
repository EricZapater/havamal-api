# Havamal API

Havamal API is a backend service built with Go (Golang) using the Gin framework. It provides authentication, user management, and content management capabilities for a blog or CMS.

## Data Models

### User

Represents a registered user of the system.

```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "is_admin": boolean,
  "is_active": boolean,
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Post

Represents a content post.

```json
{
  "id": "uuid",
  "title": "string",
  "slug": "string",
  "summary": "string",
  "content": "string",
  "status": "draft|published|archived",
  "published_at": "timestamp",
  "updated_at": "timestamp",
  "author_id": "uuid"
}
```

### Category

Represents a category for classifying posts.

```json
{
  "id": "uuid",
  "name": "string",
  "slug": "string",
  "description": "string",
  "order": int,
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Version

Represents a version history of a post.

```json
{
  "id": "uuid",
  "version": "string",
  "post_id": "uuid",
  "version_number": int,
  "content": "string",
  "created_at": "timestamp"
}
```

### Navigation

Represents a navigation menu item.

```json
{
  "id": "uuid",
  "label": "string",
  "slug": "string",
  "type": "internal|external",
  "order": int,
  "parent_id": "uuid"
}
```

## API Routes

### Authentication

Public routes for user authentication. Prefix: `/auth`

| Method | Endpoint         | Description       |
| :----- | :--------------- | :---------------- |
| `POST` | `/auth/login`    | Login user        |
| `POST` | `/auth/register` | Register new user |

### Blog API (Public)

Read-only endpoints for public consumption. Prefix: `/blog`

#### Posts

| Method | Endpoint                   | Description                |
| :----- | :------------------------- | :------------------------- |
| `GET`  | `/blog/slug/:slug`         | Get post by slug           |
| `GET`  | `/blog/author/:author_id`  | Get posts by author        |
| `GET`  | `/blog/category/:category` | Get posts by category slug |

#### Categories

| Method | Endpoint                      | Description          |
| :----- | :---------------------------- | :------------------- |
| `GET`  | `/blog/categories`            | Get all categories   |
| `GET`  | `/blog/categories/:id`        | Get category by ID   |
| `GET`  | `/blog/categories/slug/:slug` | Get category by slug |

#### Navigation

| Method | Endpoint                      | Description                 |
| :----- | :---------------------------- | :-------------------------- |
| `GET`  | `/blog/navigation`            | Get all navigation items    |
| `GET`  | `/blog/navigation/:id`        | Get navigation item by ID   |
| `GET`  | `/blog/navigation/slug/:slug` | Get navigation item by slug |

#### Versions

| Method | Endpoint             | Description       |
| :----- | :------------------- | :---------------- |
| `GET`  | `/blog/versions`     | Get all versions  |
| `GET`  | `/blog/versions/:id` | Get version by ID |

---

### Management API (Protected)

Requires Authentication (Bearer Token). Prefix: `/api`

#### Users

| Method   | Endpoint         | Description    |
| :------- | :--------------- | :------------- |
| `POST`   | `/api/users`     | Create a user  |
| `GET`    | `/api/users`     | Get all users  |
| `GET`    | `/api/users/:id` | Get user by ID |
| `PUT`    | `/api/users/:id` | Update user    |
| `DELETE` | `/api/users/:id` | Delete user    |

#### Posts

| Method   | Endpoint              | Description               |
| :------- | :-------------------- | :------------------------ |
| `POST`   | `/api/posts`          | Create a post             |
| `GET`    | `/api/posts/:id`      | Get post by ID            |
| `PUT`    | `/api/posts/:id`      | Update post               |
| `DELETE` | `/api/posts/:id`      | Delete post               |
| `POST`   | `/api/posts/category` | Add category to post      |
| `DELETE` | `/api/posts/category` | Remove category from post |
| `POST`   | `/api/posts/version`  | Add version to post       |
| `DELETE` | `/api/posts/version`  | Remove version from post  |

#### Categories

| Method   | Endpoint              | Description       |
| :------- | :-------------------- | :---------------- |
| `POST`   | `/api/categories`     | Create a category |
| `PUT`    | `/api/categories/:id` | Update category   |
| `DELETE` | `/api/categories/:id` | Delete category   |

#### Navigation

| Method   | Endpoint              | Description              |
| :------- | :-------------------- | :----------------------- |
| `POST`   | `/api/navigation`     | Create a navigation item |
| `PUT`    | `/api/navigation/:id` | Update navigation item   |
| `DELETE` | `/api/navigation/:id` | Delete navigation item   |

#### Versions

| Method   | Endpoint            | Description      |
| :------- | :------------------ | :--------------- |
| `POST`   | `/api/versions`     | Create a version |
| `PUT`    | `/api/versions/:id` | Update version   |
| `DELETE` | `/api/versions/:id` | Delete version   |

## Setup & Running

1. **Configuration**: Ensure `.env` is configured with database and auth settings.
2. **Migrations**: The application runs migrations on startup.
   - Use `cmd/migrate/main.go` for manual control: `go run cmd/migrate/main.go -action=[up|down|reset]`
3. **Run**:
   ```bash
   go run cmd/api/main.go
   ```
