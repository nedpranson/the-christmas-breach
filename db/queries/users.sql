-- name: CreateUser :one
INSERT INTO users (username, password_hash, seed)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? COLLATE NOCASE;
