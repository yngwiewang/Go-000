package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gin-gonic/gin"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("mysql",
		"root:111111@tcp(192.168.220.102:3306)/hello?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(&model.Book{})

	return db
}

func main() {
	db := initDB()
	defer db.Close()

	bookAPI := InitBookAPI(db)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/books", bookAPI.Create)
		apiv1.DELETE("/books/:id", bookAPI.Delete)
		apiv1.PUT("/books/:id", bookAPI.Update)
		apiv1.GET("/books", bookAPI.GetAll)
		apiv1.GET("/books/:id", bookAPI.GetByID)

	}

	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen err: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exiting")
}
