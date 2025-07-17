package internal

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file"})
		return
	}

	f, _ := file.Open()
	defer f.Close()

	content, _ := ioutil.ReadAll(f)
	// TODO: Extract text & embed
	_ = content

	c.JSON(http.StatusOK, gin.H{"status": "File received (stub)"})
}
