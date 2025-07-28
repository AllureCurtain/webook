package web

import (
	"github.com/gin-gonic/gin"
	"webook/pkg/ginx"
)

type Result = ginx.Result

type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type handler interface {
	RegisterRoutes(s *gin.Engine)
}
