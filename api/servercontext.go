package api

import (
	hpf "github.com/d1str0/HPFBroker"
	auth "github.com/d1str0/HPFBroker/auth"

	rbac "github.com/mikespook/gorbac"
)

type ServerContext struct {
	Version string
	*hpf.DB
	*auth.JWTSecret
	*rbac.RBAC
}
