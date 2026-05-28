CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =====================================================
-- TENANTS
-- =====================================================

CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,

    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- IDENTITIES (GLOBAL HUMAN IDENTITIES)
-- =====================================================

CREATE TABLE identities (
    id UUID PRIMARY KEY,

    primary_email TEXT NOT NULL UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT false,

    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- USERS (TENANT MEMBERSHIPS)
-- =====================================================

CREATE TABLE users (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,

    username TEXT,
    display_name TEXT,

    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (tenant_id, identity_id),
    UNIQUE (tenant_id, username)
);

-- =====================================================
-- CREDENTIALS
-- =====================================================

CREATE TABLE credentials (
    id UUID PRIMARY KEY,

    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,

    credential_type TEXT NOT NULL,

    password_hash TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- DEVICES
-- =====================================================

CREATE TABLE devices (
    id UUID PRIMARY KEY,

    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,

    device_fingerprint TEXT NOT NULL,

    platform TEXT,
    browser TEXT,
    os TEXT,

    trusted BOOLEAN NOT NULL DEFAULT false,

    last_seen_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- MFA FACTORS
-- =====================================================

CREATE TABLE mfa_factors (
    id UUID PRIMARY KEY,

    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,

    factor_type TEXT NOT NULL,

    secret_encrypted TEXT,
    public_key TEXT,

    label TEXT,

    verified BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- SESSIONS
-- =====================================================

CREATE TABLE sessions (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    refresh_token_hash TEXT NOT NULL,

    parent_session_id UUID REFERENCES sessions(id),

    ip_address INET,
    user_agent TEXT,

    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- ROLES
-- =====================================================

CREATE TABLE roles (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    name TEXT NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (tenant_id, name)
);

-- =====================================================
-- PERMISSIONS
-- =====================================================

CREATE TABLE permissions (
    id UUID PRIMARY KEY,

    name TEXT NOT NULL UNIQUE,
    description TEXT
);

-- =====================================================
-- ROLE PERMISSIONS
-- =====================================================

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,

    PRIMARY KEY (role_id, permission_id)
);

-- =====================================================
-- USER ROLES
-- =====================================================

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,

    PRIMARY KEY (user_id, role_id)
);

-- =====================================================
-- POLICIES
-- =====================================================

CREATE TABLE policies (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    name TEXT NOT NULL,

    effect TEXT NOT NULL,

    condition TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- OAUTH CLIENTS
-- =====================================================

CREATE TABLE oauth_clients (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    client_id TEXT NOT NULL UNIQUE,
    client_secret_hash TEXT,

    client_name TEXT NOT NULL,

    redirect_uris JSONB NOT NULL,

    grant_types JSONB NOT NULL,

    scopes JSONB NOT NULL,

    confidential BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- SERVICE ACCOUNTS
-- =====================================================

CREATE TABLE service_accounts (
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    name TEXT NOT NULL,

    client_id TEXT NOT NULL UNIQUE,
    client_secret_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- AUDIT EVENTS
-- =====================================================

CREATE TABLE audit_events (
    id UUID PRIMARY KEY,

    tenant_id UUID,

    actor_identity_id UUID,

    event_type TEXT NOT NULL,

    resource_type TEXT,
    resource_id UUID,

    ip_address INET,

    payload JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =====================================================
-- INDEXES
-- =====================================================

CREATE INDEX idx_users_tenant_id
ON users(tenant_id);

CREATE INDEX idx_users_identity_id
ON users(identity_id);

CREATE INDEX idx_sessions_identity_id
ON sessions(identity_id);

CREATE INDEX idx_sessions_user_id
ON sessions(user_id);

CREATE INDEX idx_sessions_expires_at
ON sessions(expires_at);

CREATE INDEX idx_sessions_revoked_at
ON sessions(revoked_at);

CREATE INDEX idx_credentials_identity_id
ON credentials(identity_id);

CREATE INDEX idx_devices_identity_id
ON devices(identity_id);

CREATE INDEX idx_mfa_factors_identity_id
ON mfa_factors(identity_id);

CREATE INDEX idx_roles_tenant_id
ON roles(tenant_id);

CREATE INDEX idx_user_roles_user_id
ON user_roles(user_id);

CREATE INDEX idx_role_permissions_role_id
ON role_permissions(role_id);

CREATE INDEX idx_audit_events_tenant_id
ON audit_events(tenant_id);

CREATE INDEX idx_audit_events_actor_identity_id
ON audit_events(actor_identity_id);

CREATE INDEX idx_audit_events_created_at
ON audit_events(created_at);