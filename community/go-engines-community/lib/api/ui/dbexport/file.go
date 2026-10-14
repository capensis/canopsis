package dbexport

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func AttachFile(c *gin.Context, collection string, content []byte) {
	c.DataFromReader(
		http.StatusOK,
		int64(len(content)),
		binding.MIMEJSON,
		bytes.NewBuffer(content),
		map[string]string{
			"Content-Disposition": "attachment; filename=" + fmt.Sprintf("%s-%d.json", collection, time.Now().Unix()),
		},
	)
}
