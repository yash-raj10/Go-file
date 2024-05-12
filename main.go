package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("ERROr")
	}

	// gin app setup
	r := gin.Default()
	r.Static("/assets", "./assets" )
	r.LoadHTMLGlob("templates/*")
	r.MaxMultipartMemory = 8 << 20

	//s3 Uploader setup
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
	return
	}
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)


	// Routes
	r.GET("/", func(c *gin.Context){
		c.HTML(http.StatusOK, "index.html", gin.H{
		})
	})

	r.POST("/", func(c *gin.Context){

		// get the file
		file, err:= c.FormFile("audio")
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error":"Failed to get the audio",
			})
			return
		} 

//-------------------------- save the file--------------------

		// save to disk--
		// err = c.SaveUploadedFile(file, "assets/uploads/"+ file.Filename)
		// if err != nil {
		// 	c.HTML(http.StatusOK, "index.html", gin.H{
		// 		"error":"Failed to Upload audio",
		// 	})
		// } 

		// save to s3-- (open file)
		f, openErr := file.Open() 
		if openErr != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error":"Failed to open audio file",
			})
			return
		} 

		// (upload file)
		result, UploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("go-file"),
			Key:    aws.String(file.Filename),
			Body:   f,
			ACL: "public-read",
		})
		if UploadErr != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error":"Failed to Upload audio",
			})
			return
		} 

		// render  the file
		c.HTML(http.StatusOK, "index.html", gin.H{
			// passing audio
			// "audio": "/assets/uploads/"+ file.Filename,
			"audio": result.Location,
		})
	})

	r.Run()
}