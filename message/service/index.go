package service

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, _ := template.ParseFiles("index.html")
	ind.Execute(c.Writer, "")
}
