CREATE TABLE users (
    id               INTEGER PRIMARY KEY,
    username         TEXT NOT NULL UNIQUE COLLATE NOCASE
                     CHECK (length(username) BETWEEN 1 AND 32),
    password_hash    TEXT NOT NULL,
    seed             INTEGER NOT NULL,
    completed_levels INTEGER NOT NULL DEFAULT 0
                     CHECK (completed_levels BETWEEN 0 AND 1023),
    created_at       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
