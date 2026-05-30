package domain

import (
	"context"
	"time"
)

type Client struct {
	ID               string
	TenantID         string
	ClientID         string
	ClientSecretHash string
	ClientName       string
	RedirectURIs     []string
	GrantTypes       []string
	Scopes           []string
	Confidential     bool
	CreatedAt        time.Time
}

func (c *Client) HasGrantType(g string) bool {
	for _, x := range c.GrantTypes {
		if x == g {
			return true
		}
	}
	return false
}

func (c *Client) HasRedirectURI(u string) bool {
	for _, x := range c.RedirectURIs {
		if x == u {
			return true
		}
	}
	return false
}

func (c *Client) AllowsScope(s string) bool {
	for _, x := range c.Scopes {
		if x == s {
			return true
		}
	}
	return false
}

type ClientRepository interface {
	GetByClientID(ctx context.Context, clientID string) (*Client, error)
	Create(ctx context.Context, client *Client) error
}
