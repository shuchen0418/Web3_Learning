package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddProduct(c *gin.Context) {

	name := c.Query("name")
	price := c.DefaultQuery("price", "0")

	c.JSON(http.StatusOK, gin.H{
		"v1":    "AddProduct",
		"name":  name,
		"price": price,
	})

}
