package common

import "github.com/gin-gonic/gin"

type CrudAPI interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type BulkCrudAPI interface {
	CrudAPI
	BulkCreate(c *gin.Context)
	BulkUpdate(c *gin.Context)
	BulkDelete(c *gin.Context)
}
