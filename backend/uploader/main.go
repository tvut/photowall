package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var FILE_ROOT string = "."

func postUpload(c *gin.Context) {

	file, handler, err := c.Request.FormFile("file")

	// Capture any errors that may arise from rest parse
	if err != nil {
		log.Println("Error getting the file")
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err})
		return
	}

	defer file.Close()

	// Printing out file information //
	log.Printf("Uploaded file name: %v\n", handler.Filename)
	log.Printf("Uploaded file size: %v\n", handler.Size)
	log.Printf("File mime type: %+v\n", handler.Header)

	// Get the file content type and access the file extension
	fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[1]
	log.Printf("File type: %v\n", fileType)
	if fileType != "jpeg" {
		log.Println("Wrong file type, exiting")
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": "Only accepts jpeg"})
		return
	}

	today := time.Now()
	filePath := fmt.Sprintf("%s/%d/%d/%d", FILE_ROOT, today.Year(), today.Month(), today.Day())

	log.Printf("Saving to: %s/%s\n", filePath, handler.Filename)
	// Done printing file information //

	// Create dir path
	err = os.MkdirAll(filePath, 0755)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	// Read file bytes
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err})
		return
	}

	// Write file to disk
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", filePath, handler.Filename), fileBytes, 0644)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	log.Printf("Successfully uploaded file")
	c.IndentedJSON(http.StatusAccepted, gin.H{"message": "success"})
}

func main() {
	router := gin.Default()
	router.POST("/upload", postUpload)

	router.Run("localhost:5000")
}
