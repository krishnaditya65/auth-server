CREATE TABLE signing_keys (
    id          TEXT PRIMARY KEY,
    algorithm   TEXT NOT NULL,
    public_key  TEXT NOT NULL,
    private_key TEXT NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    rotated_at  TIMESTAMPTZ
);

CREATE INDEX idx_signing_keys_active ON signing_keys (is_active, created_at DESC);
