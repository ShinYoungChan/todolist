package handlers

import (
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/response"
	"backend/internal/service"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTP 요청/응답 처리
type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	// 1. 데이터 바인딩 (ShouldBind는 form-data도 처리)
	var req models.TodoCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "데이터 형식이 잘못되었습니다.")
		return
	}

	var imageURL *string
	// 2. 이미지가 있다면 저장 로직 실행
	if req.Image != nil {
		// 이미지 파일 체크
		if err := checkFile(req.Image); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// 고유한 파일명 생성 (현재시간 + 원본파일명)
		newFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), req.Image.Filename)
		dst := "uploads/" + newFileName

		// 서버 로컬에 파일 저장
		if err := c.SaveUploadedFile(req.Image, dst); err != nil {
			response.Error(c, http.StatusInternalServerError, "이미지 저장 실패")
			return
		}

		// 접근 가능한 URL 경로 설정
		url := "/uploads/" + newFileName
		imageURL = &url
	}

	// 3. DB 모델로 변환
	newTodo := models.Todo{
		Title:      req.Title,
		Content:    req.Content,
		UserID:     userID,
		CategoryID: req.CategoryID,
		StartDate:  req.StartDate,
		DueDate:    req.DueDate,
		ImageURL:   imageURL, // 저장된 경로 (없으면 nil)
	}

	if err := h.service.AddTodo(&newTodo); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, gin.H{
		"message": "작성 완료",
		"todo":    newTodo,
	})
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	// 힌트: 미들웨어에서 받은 user_id로 해당 유저의 목록만 요청
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	sortBy := c.DefaultQuery("sort", "created_at")
	filter := c.DefaultQuery("filter", "all")
	keyword := c.Query("keyword")
	categoryID := c.Query("category_id")

	todos, err := h.service.GetUserTodos(userID, sortBy, filter, keyword, categoryID)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 객체로 감싸서 보내주면 유연성 증가
	response.Success(c, gin.H{"todos": todos})
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	// 힌트: URL 파라미터에서 Todo ID 추출 (c.Param("id"))
	idStr := c.Param("id")
	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID 형식이 올바르지 않습니다.")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	// 1. 수정 전 기존 데이터를 먼저 조회 (기존 이미지 경로 확인용)
	// 서비스에 FindByID 같은 함수가 구현되어 있어야 합니다.
	existingTodo, err := h.service.GetTodoByID(uint(todoID), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "항목을 찾을 수 없습니다.")
		return
	}

	// 2. form-data 바인딩 (ShouldBindJSON -> ShouldBind)
	var req models.TodoCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "잘못된 데이터 형식입니다.")
		return
	}

	// 3. 이미지 처리 로직
	var newImageURL *string = existingTodo.ImageURL // 일단 기존 경로 유지

	if req.Image != nil {
		// 이미지 파일 체크
		if err := checkFile(req.Image); err != nil {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// (1) 새 이미지 저장
		newFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), req.Image.Filename)
		dst := "uploads/" + newFileName
		if err := c.SaveUploadedFile(req.Image, dst); err == nil {
			// (2) 새 저장 성공 시, 기존 이미지가 있었다면 삭제
			if existingTodo.ImageURL != nil {
				deleteFile(*existingTodo.ImageURL)
			}
			// (3) 경로 업데이트
			url := "/uploads/" + newFileName
			newImageURL = &url
		}
	}

	// 4. 기존 Todo 객체에 새 데이터 덮어쓰기
	existingTodo.Title = req.Title
	existingTodo.Content = req.Content
	existingTodo.CategoryID = req.CategoryID
	existingTodo.StartDate = req.StartDate
	existingTodo.DueDate = req.DueDate
	existingTodo.ImageURL = newImageURL

	// 5. DB 업데이트 호출
	// 서비스 레이어의 UpdateTodo 함수가 existingTodo 자체를 받도록 수정되어 있다면 더 편합니다.
	if err := h.service.UpdateTodo(uint(todoID), userID, existingTodo); err != nil {
		response.Error(c, http.StatusInternalServerError, "수정 실패")
		return
	}

	response.SuccessMsg(c, "변경 성공")
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	// 힌트: 삭제할 권한(내 글이 맞는지) 확인 로직 고려
	idStr := c.Param("id")
	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID 형식이 올바르지 않습니다.")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	// 1. 삭제 전 데이터 조회 (이미지 경로를 확인하기 위해)
	todo, err := h.service.GetTodoByID(uint(todoID), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "삭제할 항목을 찾을 수 없습니다.")
		return
	}

	// 2. 이미지가 있다면 서버 로컬 폴더에서 실제 파일 삭제
	if todo.ImageURL != nil {
		deleteFile(*todo.ImageURL)
	}

	if err := h.service.DeleteTodo(uint(todoID), userID); err != nil {
		response.Error(c, http.StatusForbidden, "삭제 권한이 없거나 항목을 찾을 수 없습니다.")
		return
	}

	response.SuccessMsg(c, "삭제 성공")
}

func deleteFile(filePath string) {
	if filePath == "" {
		return
	}
	// 상대 경로이므로 앞에 .을 붙이거나 저장된 경로 그대로 사용 (저장 방식에 따라 조절)
	// "/uploads/xxx.jpg" -> "uploads/xxx.jpg" 로 변경이 필요할 수 있습니다.
	actualPath := strings.TrimPrefix(filePath, "/")

	if err := os.Remove(actualPath); err != nil {
		log.Printf("파일 삭제 실패: %v", err)
	}
}

func checkFile(file *multipart.FileHeader) error {
	// 1. 용량 체크 (타입 캐스팅 방지를 위해 처음부터 int64로 선언)
	const maxFileSize = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		return errors.New("파일 용량은 5MB를 초과할 수 없습니다.")
	}

	// 2. 확장자 체크 (Map을 사용하면 반복문을 돌지 않아도 되어 더 빠르게 동작가능)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return errors.New("파일 확장자가 잘못되었습니다.")
	}

	// 3. 마임 타입 체크
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return errors.New("파일 타입이 잘못되었습니다.")
	}

	return nil
}
