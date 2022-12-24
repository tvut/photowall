package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var FILE_ROOT string = "."

func postUpload(c *gin.Context) {

	file, handler, err := c.Request.FormFile("file")

	// Capture any errors that may arise
	if err != nil {
		print("Error getting the file")
		print(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded file name: %v\n", handler.Filename)
	fmt.Printf("Uploaded file size: %v\n", handler.Size)
	fmt.Printf("File mime type: %+v\n", handler.Header)

	// Get the file content type and access the file extension
	fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[1]
	fmt.Printf("File type: %v\n", fileType)

	today := time.Now()
	filePath := fmt.Sprintf("%s/%d/%d/%d", FILE_ROOT, today.Year(), today.Month(), today.Day())

	fmt.Printf("Saving to: %s/%s\n", filePath, handler.Filename)

	err = os.MkdirAll(filePath, 0755)
	if err != nil {
		log.Println(err)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", filePath, handler.Filename), fileBytes, 0644)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Successfully uploaded file")
}

func main() {
	router := gin.Default()
	router.POST("/upload", postUpload)

	router.Run("localhost:5000")
}
