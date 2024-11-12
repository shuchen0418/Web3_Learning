package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Member struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

func AddMember(c *gin.Context) {
	var member Member

	member.Name = c.Query("name")
	member.Price = c.DefaultQuery("price", "0")

	c.JSON(http.StatusOK, member)
}
