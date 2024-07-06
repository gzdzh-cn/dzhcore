package dzhcore

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	IsRefresh       bool     `json:"isRefresh"`
	RoleIds         []string `json:"roleIds"`
	Username        string   `json:"username"`
	UserId          string   `json:"userId"`
	PasswordVersion *int32   `json:"passwordVersion"`
	jwt.RegisteredClaims
}
