package hpfbroker

import (
	auth "github.com/d1str0/HPFBroker/auth"
	rbac "github.com/mikespook/gorbac"
)

type ServerContext struct {
	Version string
	*DB
	*auth.JWTSecret
	*rbac.RBAC
}
