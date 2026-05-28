package tenancy

import "context"

type contextKey string

const tenantKey contextKey = "tenant_id"

func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

func TenantID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(tenantKey).(string)
	return v, ok
}
