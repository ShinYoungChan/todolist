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

	if err := db.AutoMigrate(&models.User{}, &models.Category{}, &models.Todo{}, &models.RefreshToken{}); err != nil {
		log.Fatalf("마이그레이션 실패: %v", err)
	}

	// 의존성 주입
	tokenRepo := repository.NewTokenRepository(db)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, tokenRepo, jwtSecret)
	userHandler := handlers.NewUserHandler(userService)

	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	r := gin.Default()

	// 💡 CORS 미들웨어 추가 (여기가 핵심!)
	r.Use(func(c *gin.Context) {
		// 1. 모든 도메인 허용 (Flutter Web 개발 시 필수)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 2. PUT을 포함한 모든 메서드 허용
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// 3. Authorization 헤더 허용 (JWT 토큰 사용 시 필수)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		// 4. 자격 증명 허용
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 브라우저의 OPTIONS 요청(Preflight)에 대해 200 또는 204로 즉시 응답
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200) // 204 대신 200을 써도 무방합니다.
			return
		}

		c.Next()
	})

	routes.SetupUserRoutes(r, userHandler, jwtSecret)
	routes.SetupTodoRoutes(r, todoHandler, jwtSecret)
	routes.SetupCategoryRoutes(r, categoryHandler, jwtSecret)

	r.Run(":" + port)
}
