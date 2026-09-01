CREATE TABLE tokens (
    hash       BLOB PRIMARY KEY,
    scope      TEXT NOT NULL CHECK (scope IN ('activation', 'password_reset', 'authentication')),
    user_id    TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ'))
) STRICT;
