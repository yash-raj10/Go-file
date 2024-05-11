package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main(){
	r := gin.Default()
	r.Static("/assets", "./assets" )
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context){
		c.HTML(http.StatusOK, "index.html", gin.H{
		})
	})

	r.POST("/", func(c *gin.Context){

		// get the file
		file, err:= c.FormFile("audio")
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error":"Failed to Upload audio",
			})
		} 

		// sve the file
		err = c.SaveUploadedFile(file, "assets/uploads/"+ file.Filename)
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error":"Failed to Upload audio",
			})
		} 

		// render  the file
		c.HTML(http.StatusOK, "index.html", gin.H{
			// passing audio
			"audio":"/assets/uploads/"+ file.Filename,
		})
	})

	r.Run()
}