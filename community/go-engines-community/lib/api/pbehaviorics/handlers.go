package pbehaviorics

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GetICS(store Store, service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		pbh, err := store.GetOneBy(c.Request.Context(), id)
		if err != nil {
			panic(err)
		}

		maxPriority, err := store.FindMaxPriority(c.Request.Context())
		if err != nil {
			panic(err)
		}

		minPriority, err := store.FindMinPriority(c.Request.Context())
		if err != nil {
			panic(err)
		}

		if pbh == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		}

		calendar, err := service.GenICSFrom(*pbh, maxPriority, minPriority)
		if err != nil {
			panic(err)
		}

		fileName := genFileName(pbh.Name)
		extraHeaders := map[string]string{
			"Content-Disposition": "attachment; filename=" + fileName,
		}
		contentType := "text/calendar"
		reader, contentLength := calendar.Read()

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	}
}

func genFileName(pbehaviorName string) string {
	filename := pbehaviorName
	// Normalize unicode chars
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	filename, _, _ = transform.String(t, filename)
	// Transform to lower case
	filename = strings.ToLower(filename)
	// Replace all whitspace to "-"
	r := regexp.MustCompile(`[\s]+`)
	filename = r.ReplaceAllString(filename, "-")
	// Remove all not alphanumeric chars
	r = regexp.MustCompile(`[^a-z0-9\-]+`)
	filename = r.ReplaceAllString(filename, "")

	filename = fmt.Sprintf("pbehavior-%s.ics", filename)

	return filename
}
