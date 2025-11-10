package crud

import "github.com/gin-gonic/gin"

type API interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type BulkAPI interface {
	API
	BulkCreate(c *gin.Context)
	BulkUpdate(c *gin.Context)
	BulkDelete(c *gin.Context)
}
