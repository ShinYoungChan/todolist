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
	if err := godotenv.Load("go.env"); err != nil {
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET이 설정되지 않았습니다!")
	}

	// 3. DB 연결 및 마이그레이션
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		log.Fatalf("마이그레이션 실패: %v", err)
	}

	// 의존성 주입
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtSecret)
	userHandler := handlers.NewUserHandler(userService)

	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	r := gin.Default()

	// 💡 CORS 미들웨어 추가 (여기가 핵심!)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routes.SetupUserRoutes(r, userHandler)
	routes.SetupTodoRoutes(r, todoHandler, jwtSecret)

	r.Run(":" + port)
}
