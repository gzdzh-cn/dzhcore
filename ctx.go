package dzhcore

import (
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ctx = gctx.GetInitCtx()
)

type Claims struct {
	IsRefresh       bool     `json:"isRefresh"`
	RoleIds         []string `json:"roleIds"`
	Username        string   `json:"username"`
	UserId          string   `json:"userId"`
	PasswordVersion *int32   `json:"passwordVersion"`
	jwt.RegisteredClaims
}
