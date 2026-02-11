package main

import (
	"backend/internal/handlers"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/internal/service"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Main Start")

	// 1. .env 로드 (파일이 없어도 서버는 뜰 수 있게Fatal 대신 Print)
	if err := godotenv.Load(); err != nil {
		log.Println("경고: .env 파일을 찾을 수 없습니다. 기본 설정을 사용합니다.")
	}

	// 2. 환경 변수 읽기
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "todo.db"
	}

	// 3. DB 연결 및 마이그레이션
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("마이그레이션 실패: %v", err)
	}

	// 의존성 주입
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	routes.SetupUserRoutes(r, userHandler)

	r.Run(":" + port)
}
