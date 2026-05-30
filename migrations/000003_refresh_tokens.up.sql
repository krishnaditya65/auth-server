CREATE TABLE refresh_tokens (
    id TEXT PRIMARY KEY,

    session_id TEXT NOT NULL
        REFERENCES sessions(id),

    token_hash TEXT NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,

    revoked_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_refresh_tokens_session_id
ON refresh_tokens(session_id);

CREATE INDEX idx_refresh_tokens_token_hash
ON refresh_tokens(token_hash);