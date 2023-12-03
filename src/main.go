package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	log.Infof("Starting...")

	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		log.Fatalf("Failed to retrieve MySQL Host from Environment.")
	}
	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlHost == "" {
		log.Fatalf("Failed to retrieve MySQL User from Environment.")
	}
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlHost == "" {
		log.Fatalf("Failed to retrieve MySQL Password from Environment.")
	}
	mysqlDbName := os.Getenv("MYSQL_DBNAME")
	if mysqlHost == "" {
		log.Fatalf("Failed to retrieve MySQL DB Name from Environment.")
	}

	log.Infof("Configuration Loaded.")

	// Construct MySQL connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlDbName)

	// Initialize and connect to the MySQL database
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate schema: %v", err)
	}

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(config))
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8080")
}
