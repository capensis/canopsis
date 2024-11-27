package dbexport

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func AttachFile(c *gin.Context, collection string, content []byte) (err error) {
	file, err := os.CreateTemp("", collection+"_*")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			err = fmt.Errorf("failed to close the file %s: %w", file.Name(), err)
			return
		}

		err = os.Remove(file.Name())
		if err != nil {
			err = fmt.Errorf("failed to remove the file %s: %w", file.Name(), err)
		}
	}()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write the file %s: %w", file.Name(), err)
	}

	c.FileAttachment(file.Name(), fmt.Sprintf("%s-%d.json", collection, time.Now().Unix()))

	return nil
}
