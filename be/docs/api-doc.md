# API Documentation

**Base URL:** `/api`  
**Authentication:** JWT Bearer Token (except public endpoints)  
**Content-Type:** `application/json` (default), `multipart/form-data` (for file uploads)

## Authentication

Include the JWT token in the `Authorization` header for protected routes:

```
Authorization: Bearer <token>
```

## Response Format

### Success Response

```json
{
  "code": 200,
  "message": "Success message",
  "data": { ... }
}
```

### Paginated Response

```json
{
  "code": 200,
  "message": "Success message",
  "data": [ ... ],
  "metadata": {
    "total": 50,
    "page": 1,
    "page_size": 10,
    "total_pages": 5
  }
}
```

### Error Response

```json
{
  "code": 400,
  "message": "Error message"
}
```

### Validation Error Response (422)

```json
{
  "code": 422,
  "message": "The given data was invalid.",
  "errors": {
    "email": ["Email is required"],
    "name": ["Name must be at least 2 characters"]
  }
}
```

---

## System Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | Public | Health check |
| GET | `/.well-known/jwks.json` | Public | Get JSON Web Key Set |
| POST | `/api/upload` | JWT | Upload file |

### GET `/health`

Check the health status of the application and its dependencies.

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Health check successful",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T00:00:00Z",
    "database": {
      "status": "healthy",
      "latency": "12ms"
    },
    "redis": {
      "status": "healthy",
      "latency": "5ms"
    }
  }
}
```

---

### GET `/.well-known/jwks.json`

Retrieve the JSON Web Key Set for token verification.

**Response (200 OK)**

```json
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "...",
      "use": "sig",
      "alg": "RS256",
      "n": "...",
      "e": "..."
    }
  ]
}
```

---

### POST `/api/upload`

Upload a file to temporary storage. The returned `uuid` can be used in user create/update endpoints via `avatar`.

**Headers**

```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request Body** (multipart/form-data)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `file` | file | Yes | jpg, jpeg, png, gif, webp, jfif. Max 5MB |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "File uploaded successfully",
  "data": {
    "uuid": "550e8400-e29b-41d4-a716-446655440000",
    "url": "http://localhost:8080/storage/tmp/550e8400-e29b-41d4-a716-446655440000.jpg"
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | File is required, file type not allowed, or file size exceeds limit |
| 401 | Unauthorized |
| 500 | Internal server error |

---

## Auth Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/auth/login` | Public | Login user |
| POST | `/api/auth/refresh` | Public | Refresh access token |
| POST | `/api/auth/forgot-password` | Public | Send password reset email |
| POST | `/api/auth/reset-password` | Public | Reset password with token |
| POST | `/api/auth/logout` | JWT | Logout user |

### POST `/api/auth/login`

Authenticate user and receive access tokens.

**Request Body**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `email` | string | Yes | Valid email format |
| `password` | string | Yes | Min 6 characters |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "avatar": "http://localhost:8080/storage/avatars/550e8400-e29b-41d4-a716-446655440000.jpg",
      "created_at": "2024-01-01T00:00:00Z",
      "roles": [
        {
          "id": 1,
          "name": "admin",
          "description": "Administrator role"
        }
      ],
      "permissions": [
        {
          "id": 1,
          "name": "users.read",
          "description": "View users"
        }
      ]
    }
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid credentials |
| 422 | Validation error (email/password format) |
| 500 | Internal server error |

---

### POST `/api/auth/refresh`

Refresh an expired access token using a refresh token.

**Request Body**

```json
{
  "refresh_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

| Field | Type | Required |
|-------|------|----------|
| `refresh_token` | string | Yes |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJSUzI1NiIs..."
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid or expired refresh token |
| 422 | Validation error |
| 500 | Internal server error |

---

### POST `/api/auth/forgot-password`

Request a password reset email.

**Request Body**

```json
{
  "email": "user@example.com"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `email` | string | Yes | Valid email format |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Reset password email sent",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 404 | Email not found |
| 422 | Validation error |
| 500 | Internal server error |

---

### POST `/api/auth/reset-password`

Reset password using the token received via email.

**Request Body**

```json
{
  "token": "reset-token-from-email",
  "new_password": "newpassword123"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `token` | string | Yes | - |
| `new_password` | string | Yes | Min 6 characters |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Password reset successful",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid or expired token |
| 422 | Validation error |
| 500 | Internal server error |

---

### POST `/api/auth/logout`

Logout the current user.

**Headers**

```
Authorization: Bearer <token>
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Logout successful. Please clear your token on the client side.",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 401 | Unauthorized (missing/invalid token) |
| 500 | Internal server error |

---

## Permissions Endpoints

| Method | Path | Auth | Permission Required | Description |
|--------|------|------|---------------------|-------------|
| POST | `/api/permissions` | JWT | `permission.create` | Create permission |
| GET | `/api/permissions` | JWT | `permission.index` | List permissions (all or paginated) |
| GET | `/api/permissions/:id` | JWT | `permission.index` | Get permission by ID |
| PUT | `/api/permissions/:id` | JWT | `permission.edit` | Update permission |
| DELETE | `/api/permissions/:id` | JWT | `permission.delete` | Delete permission |

### POST `/api/permissions`

Create a new permission.

**Request Body**

```json
{
  "name": "users.create",
  "description": "Create new users"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 3-100 characters |
| `description` | string | No | Max 255 characters |

**Response (201 Created)**

```json
{
  "code": 201,
  "message": "Permission created successfully",
  "data": {
    "id": 1,
    "name": "users.create",
    "description": "Create new users"
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 403 | Forbidden (missing permission) |
| 422 | Validation error |
| 500 | Internal server error |

---

### GET `/api/permissions`

List permissions. Supports optional pagination.

**Query Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number (triggers pagination) |
| `page_size` | integer | No | Items per page (default: 10) |

**Without pagination** (no `page` parameter):

```
GET /api/permissions
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Permissions fetched successfully",
  "data": [
    {
      "id": 1,
      "name": "users.create",
      "description": "Create new users"
    },
    {
      "id": 2,
      "name": "users.read",
      "description": "View users"
    }
  ]
}
```

**With pagination** (with `page` parameter):

```
GET /api/permissions?page=1&page_size=10
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Permissions fetched successfully",
  "data": [
    {
      "id": 1,
      "name": "users.create",
      "description": "Create new users"
    }
  ],
  "metadata": {
    "total": 25,
    "page": 1,
    "page_size": 10,
    "total_pages": 3
  }
}
```

---

### GET `/api/permissions/:id`

Get a single permission by ID.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Permission ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Permission fetched successfully",
  "data": {
    "id": 1,
    "name": "users.create",
    "description": "Create new users"
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid permission ID |
| 404 | Permission not found |
| 500 | Internal server error |

---

### PUT `/api/permissions/:id`

Update an existing permission.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Permission ID |

**Request Body**

```json
{
  "name": "users.create",
  "description": "Create new users"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 3-100 characters |
| `description` | string | No | Max 255 characters |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Permission updated successfully",
  "data": {
    "id": 1,
    "name": "users.create",
    "description": "Create new users"
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid permission ID |
| 404 | Permission not found |
| 422 | Validation error |
| 500 | Internal server error |

---

### DELETE `/api/permissions/:id`

Soft delete a permission.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Permission ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Permission deleted successfully",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid permission ID |
| 404 | Permission not found |
| 500 | Internal server error |

---

## Roles Endpoints

| Method | Path | Auth | Permission Required | Description |
|--------|------|------|---------------------|-------------|
| POST | `/api/roles` | JWT | `role.create` | Create role |
| GET | `/api/roles` | JWT | `role.index` | List roles (all or paginated) |
| GET | `/api/roles/:id` | JWT | `role.index` | Get role by ID |
| PUT | `/api/roles/:id` | JWT | `role.edit` | Update role |
| DELETE | `/api/roles/:id` | JWT | `role.delete` | Delete role |
| GET | `/api/roles/:id/permissions` | JWT | `role.index` | Get role permissions |

### POST `/api/roles`

Create a new role with permissions.

**Request Body**

```json
{
  "name": "admin",
  "description": "Administrator role",
  "permissions": [1, 2, 3]
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 3-100 characters |
| `description` | string | No | Max 255 characters |
| `permissions` | array | Yes | Array of permission IDs |

**Response (201 Created)**

```json
{
  "code": 201,
  "message": "Role created successfully",
  "data": {
    "id": 1,
    "name": "admin",
    "description": "Administrator role",
    "permissions": [
      {
        "id": 1,
        "name": "users.create",
        "description": "Create new users"
      },
      {
        "id": 2,
        "name": "users.read",
        "description": "View users"
      }
    ]
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 403 | Forbidden (missing permission) |
| 422 | Validation error |
| 500 | Internal server error |

---

### GET `/api/roles`

List roles. Supports optional pagination.

**Query Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number (triggers pagination) |
| `page_size` | integer | No | Items per page (default: 10) |

**Without pagination** (no `page` parameter):

```
GET /api/roles
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Roles fetched successfully",
  "data": [
    {
      "id": 1,
      "name": "admin",
      "description": "Administrator role",
      "permissions": [
        {
          "id": 1,
          "name": "users.create",
          "description": "Create new users"
        }
      ]
    },
    {
      "id": 2,
      "name": "user",
      "description": "Regular user role",
      "permissions": [
        {
          "id": 2,
          "name": "users.read",
          "description": "View users"
        }
      ]
    }
  ]
}
```

**With pagination** (with `page` parameter):

```
GET /api/roles?page=1&page_size=10
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Roles fetched successfully",
  "data": [
    {
      "id": 1,
      "name": "admin",
      "description": "Administrator role",
      "permissions": [
        {
          "id": 1,
          "name": "users.create",
          "description": "Create new users"
        }
      ]
    }
  ],
  "metadata": {
    "total": 15,
    "page": 1,
    "page_size": 10,
    "total_pages": 2
  }
}
```

---

### GET `/api/roles/:id`

Get a single role by ID.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Role ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Role fetched successfully",
  "data": {
    "id": 1,
    "name": "admin",
    "description": "Administrator role",
    "permissions": [
      {
        "id": 1,
        "name": "users.create",
        "description": "Create new users"
      }
    ]
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid role ID |
| 404 | Role not found |
| 500 | Internal server error |

---

### PUT `/api/roles/:id`

Update an existing role with permissions.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Role ID |

**Request Body**

```json
{
  "name": "super-admin",
  "description": "Super Administrator role",
  "permissions": [1, 2, 3, 4]
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 3-100 characters |
| `description` | string | No | Max 255 characters |
| `permissions` | array | Yes | Array of permission IDs |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Role updated successfully",
  "data": {
    "id": 1,
    "name": "super-admin",
    "description": "Super Administrator role",
    "permissions": [
      {
        "id": 1,
        "name": "users.create",
        "description": "Create new users"
      }
    ]
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid role ID |
| 404 | Role not found |
| 422 | Validation error |
| 500 | Internal server error |

---

### DELETE `/api/roles/:id`

Soft delete a role (clears permissions association first).

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Role ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Role deleted successfully",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid role ID |
| 404 | Role not found |
| 500 | Internal server error |

---

### GET `/api/roles/:id/permissions`

Get all permissions assigned to a role.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Role ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Role permissions fetched successfully",
  "data": [
    {
      "id": 1,
      "name": "users.create",
      "description": "Create new users"
    },
    {
      "id": 2,
      "name": "users.read",
      "description": "View users"
    }
  ]
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid role ID |
| 500 | Internal server error |

---

## Users Endpoints

| Method | Path | Auth | Permission Required | Description |
|--------|------|------|---------------------|-------------|
| POST | `/api/users` | JWT | `user.create` | Create user |
| GET | `/api/users` | JWT | `user.index` | List users (all or paginated) |
| GET | `/api/users/:id` | JWT | `user.index` | Get user by ID |
| PUT | `/api/users/:id` | JWT | `user.edit` | Update user |
| DELETE | `/api/users/:id` | JWT | `user.delete` | Delete user |

**Note:** User endpoints use JSON body. To set an avatar, first upload a file via `POST /api/upload` and use the returned `uuid` as `avatar`.

---

## Profile Endpoints

| Method | Path | Auth | Permission Required | Description |
|--------|------|------|---------------------|-------------|
| GET | `/api/me` | JWT | None | Get current user profile |
| PUT | `/api/me` | JWT | None | Update current user profile |

**Note:** Profile endpoints allow authenticated users to manage their own profile without requiring specific permissions. To set an avatar, first upload a file via `POST /api/upload` and use the returned `uuid` as `avatar`.

### GET `/api/me`

Get the authenticated user's profile information.

**Headers**

```
Authorization: Bearer <token>
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Profile fetched successfully",
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe",
    "avatar": "http://localhost:8080/storage/avatars/550e8400-e29b-41d4-a716-446655440000.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      }
    ],
    "permissions": [
      {
        "id": 1,
        "name": "users.create",
        "description": "Create new users"
      }
    ]
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 401 | Unauthorized (missing/invalid token) |
| 404 | Profile not found |
| 500 | Internal server error |

---

### PUT `/api/me`

Update the authenticated user's profile information.

**Headers**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body**

```json
{
  "name": "John Doe Updated",
  "password": "newpassword123",
  "password_confirmation": "newpassword123",
  "avatar": "660e8400-e29b-41d4-a716-446655440001"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 2-100 characters |
| `password` | string | No | Min 6 characters (empty = no change) |
| `password_confirmation` | string | No | Must match password if provided |
| `avatar` | string | No | UUID from upload endpoint |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe Updated",
    "avatar": "http://localhost:8080/storage/avatars/660e8400-e29b-41d4-a716-446655440001.png",
    "created_at": "2024-01-01T00:00:00Z",
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      }
    ],
    "permissions": [
      {
        "id": 1,
        "name": "users.create",
        "description": "Create new users"
      }
    ]
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 401 | Unauthorized (missing/invalid token) |
| 404 | Profile not found |
| 422 | Validation error |
| 500 | Internal server error |

---

### POST `/api/users`

Create a new user with roles.

**Headers**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "password_confirmation": "password123",
  "avatar": "550e8400-e29b-41d4-a716-446655440000",
  "roles": [1, 2]
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 2-100 characters |
| `email` | string | Yes | Valid email format |
| `password` | string | Yes | Min 6 characters |
| `password_confirmation` | string | Yes | Must match password |
| `avatar` | string | No | UUID from upload endpoint |
| `roles` | array | No | Array of role IDs |

**Response (201 Created)**

```json
{
  "code": 201,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe",
    "avatar": "http://localhost:8080/storage/avatars/550e8400-e29b-41d4-a716-446655440000.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      }
    ],
    "permissions": [
      {
        "id": 1,
        "name": "users.create",
        "description": "Create new users"
      }
    ]
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 403 | Forbidden (missing permission) |
| 422 | Validation error |
| 500 | Internal server error |

---

### GET `/api/users`

List users. Supports optional pagination.

**Query Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number (triggers pagination) |
| `page_size` | integer | No | Items per page (default: 10) |

**Without pagination** (no `page` parameter):

```
GET /api/users
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Users fetched successfully",
  "data": [
    {
      "id": 1,
      "email": "john@example.com",
      "name": "John Doe",
      "avatar": "http://localhost:8080/storage/avatars/550e8400-e29b-41d4-a716-446655440000.jpg",
      "created_at": "2024-01-01T00:00:00Z",
      "roles": [
        {
          "id": 1,
          "name": "admin",
          "description": "Administrator role"
        }
      ],
      "permissions": [
        {
          "id": 1,
          "name": "users.create",
          "description": "Create new users"
        }
      ]
    }
  ]
}
```

**With pagination** (with `page` parameter):

```
GET /api/users?page=1&page_size=10
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Users fetched successfully",
  "data": [
    {
      "id": 1,
      "email": "john@example.com",
      "name": "John Doe",
      "avatar": "http://localhost:8080/storage/avatars/550e8400-e29b-41d4-a716-446655440000.jpg",
      "created_at": "2024-01-01T00:00:00Z",
      "roles": [],
      "permissions": []
    }
  ],
  "metadata": {
    "total": 100,
    "page": 1,
    "page_size": 10,
    "total_pages": 10
  }
}
```

---

### GET `/api/users/:id`

Get a single user by ID.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | User ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "User fetched successfully",
  "data": {
    "id": 1,
    "email": "john@example.com",
    "name": "John Doe",
    "avatar": "http://localhost:8080/storage/avatars/550e8400-e29b-41d4-a716-446655440000.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      }
    ],
    "permissions": [
      {
        "id": 1,
        "name": "users.create",
        "description": "Create new users"
      }
    ]
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid user ID |
| 404 | User not found |
| 500 | Internal server error |

---

### PUT `/api/users/:id`

Update an existing user with roles.

**Headers**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | User ID |

**Request Body**

```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com",
  "password": "newpassword123",
  "password_confirmation": "newpassword123",
  "avatar": "660e8400-e29b-41d4-a716-446655440001",
  "roles": [1, 3]
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 2-100 characters |
| `email` | string | Yes | Valid email format |
| `password` | string | No | Min 6 characters (empty = no change) |
| `password_confirmation` | string | No | Must match password if provided |
| `avatar` | string | No | UUID from upload endpoint |
| `roles` | array | No | Array of role IDs (replaces all existing roles) |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "email": "john.updated@example.com",
    "name": "John Doe Updated",
    "avatar": "http://localhost:8080/storage/avatars/660e8400-e29b-41d4-a716-446655440001.png",
    "created_at": "2024-01-01T00:00:00Z",
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      }
    ],
    "permissions": [
      {
        "id": 1,
        "name": "users.create",
        "description": "Create new users"
      }
    ]
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid user ID |
| 404 | User not found |
| 422 | Validation error |
| 500 | Internal server error |

---

### DELETE `/api/users/:id`

Soft delete a user.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | User ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "User deleted successfully",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid user ID |
| 404 | User not found |
| 500 | Internal server error |

---

## Notifications Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/notifications` | JWT | List notifications (paginated, default page 1) |
| GET | `/api/notifications/unread-count` | JWT | Get unread notification count |
| GET | `/api/notifications/:id` | JWT | Get notification by ID |
| PATCH | `/api/notifications/:id/read` | JWT | Mark notification as read |
| PATCH | `/api/notifications/mark-all-read` | JWT | Mark all notifications as read |
| DELETE | `/api/notifications/:id` | JWT | Delete notification |

**Note:** All notification endpoints return only the authenticated user's own notifications. No additional permission required beyond authentication.

### GET `/api/notifications`

List notifications for the authenticated user. Always returns paginated results.

**Query Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | integer | No | Page number (default: 1) |
| `page_size` | integer | No | Items per page (default: 10) |
| `is_read` | string | No | Filter by read status: `true` or `false` |
| `type` | string | No | Filter by notification type (e.g., `order_completed`) |

**Example Request**

```
GET /api/notifications?page=1&page_size=10&is_read=false&type=order_completed
```

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Notifications fetched successfully",
  "data": [
    {
      "id": 1,
      "user_id": 5,
      "type": "order_completed",
      "title": "Order Completed",
      "message": "Your order #123 has been completed",
      "data": "{\"order_id\":123,\"total\":99.99}",
      "read_at": null,
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "user_id": 5,
      "type": "password_changed",
      "title": "Password Changed",
      "message": "Your password has been changed successfully",
      "data": null,
      "read_at": "2024-01-15T11:00:00Z",
      "created_at": "2024-01-14T09:00:00Z"
    }
  ],
  "metadata": {
    "total": 25,
    "page": 1,
    "page_size": 10,
    "total_pages": 3
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 401 | Unauthorized (missing/invalid token) |
| 500 | Internal server error |

---

### GET `/api/notifications/unread-count`

Get the count of unread notifications for the authenticated user.

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Unread count fetched successfully",
  "data": {
    "unread_count": 5
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 401 | Unauthorized (missing/invalid token) |
| 500 | Internal server error |

---

### GET `/api/notifications/:id`

Get a single notification by ID. Only returns the notification if it belongs to the authenticated user.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Notification ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Notification fetched successfully",
  "data": {
    "id": 1,
    "user_id": 5,
    "type": "order_completed",
    "title": "Order Completed",
    "message": "Your order #123 has been completed",
    "data": "{\"order_id\":123,\"total\":99.99}",
    "read_at": null,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid notification ID |
| 401 | Unauthorized (missing/invalid token) |
| 403 | You don't have access to this notification |
| 404 | Notification not found |
| 500 | Internal server error |

---

### PATCH `/api/notifications/:id/read`

Mark a single notification as read. Only works for notifications belonging to the authenticated user.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Notification ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Notification marked as read",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid notification ID |
| 401 | Unauthorized (missing/invalid token) |
| 404 | Notification not found |
| 500 | Internal server error |

---

### PATCH `/api/notifications/mark-all-read`

Mark all unread notifications as read for the authenticated user.

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "All notifications marked as read",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 401 | Unauthorized (missing/invalid token) |
| 500 | Internal server error |

---

### DELETE `/api/notifications/:id`

Soft delete a notification. Only works for notifications belonging to the authenticated user.

**Path Parameters**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | Yes | Notification ID |

**Response (200 OK)**

```json
{
  "code": 200,
  "message": "Notification deleted successfully",
  "data": null
}
```

**Error Responses**

| Status | Message |
|--------|---------|
| 400 | Invalid notification ID |
| 401 | Unauthorized (missing/invalid token) |
| 403 | You don't have access to this notification |
| 404 | Notification not found |
| 500 | Internal server error |

---

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 422 | Unprocessable Entity (Validation Error) |
| 500 | Internal Server Error |
