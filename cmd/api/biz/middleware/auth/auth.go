package auth

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/hertz-contrib/jwt"
)

type Middleware struct {
	JWT *jwt.HertzJWTMiddleware
}

func NewMiddleware(middleware *jwt.HertzJWTMiddleware) *Middleware {
	authMiddleware, err := jwt.New(middleware)
	if err != nil {
		logger.Fatal(err)
	}
	return &Middleware{
		JWT: authMiddleware,
	}
}
