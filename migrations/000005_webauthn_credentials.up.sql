CREATE TABLE webauthn_credentials (
    id              UUID PRIMARY KEY,
    identity_id     UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    credential_id   BYTEA NOT NULL UNIQUE,
    public_key      BYTEA NOT NULL,
    attestation     TEXT,
    aaguid          BYTEA,
    sign_count      BIGINT NOT NULL DEFAULT 0,
    transports      JSONB,
    user_handle     BYTEA NOT NULL,
    label           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at    TIMESTAMPTZ
);

CREATE INDEX idx_webauthn_credentials_identity ON webauthn_credentials (identity_id);
CREATE INDEX idx_webauthn_credentials_user_handle ON webauthn_credentials (user_handle);
