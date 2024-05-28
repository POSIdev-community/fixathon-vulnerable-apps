package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"phdays-app/src/api"
	"phdays-app/src/models"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSearchPage_Success(t *testing.T) {
	r := setUpRouter(true)
	r.GET("/search", api.Search)
	req, _ := http.NewRequest("GET", "/search", nil)
	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		pageOK := strings.Index(w.Body.String(), "<title>Articles Search</title>") > 0

		return statusOK && pageOK
	})
}

func TestSearchArticles_Success_Row(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs("%"+mockArticle.Title+"%", "%"+mockArticle.Title+"%").
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}).
			AddRow(mockArticle.ArticleId, mockArticle.Content, mockArticle.Title, mockArticle.UserId, mockArticle.Author))

	// Set up the Gin router to handle the request
	apiGroup := r.Group("/api")
	apiGroup.POST("/search", api.SearchArticles)
	// Create a new HTTP request with form data
	payload := strings.NewReader(`{"search":"title"}`)
	req, _ := http.NewRequest("POST", "/api/search", payload)
	req.Header.Set("Content-Type", "application/json")
	// req.Body = strings.NewReader(`{"search":"title"}`)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		var articles []models.Article
		err := json.NewDecoder(w.Body).Decode(&articles)
		validData := articles[0] == mockArticle
		return statusOK && err == nil && validData
	})
}

func TestSearchArticles_Success_No_Rows(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs("%"+mockArticle.Title+"%", "%"+mockArticle.Title+"%").
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}))

	// Set up the Gin router to handle the request
	apiGroup := r.Group("/api")
	apiGroup.POST("/search", api.SearchArticles)
	// Create a new HTTP request with form data
	payload := strings.NewReader(`{"search":"title"}`)
	req, _ := http.NewRequest("POST", "/api/search", payload)
	req.Header.Set("Content-Type", "application/json")
	// req.Body = strings.NewReader(`{"search":"title"}`)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		var articles []models.Article
		err := json.NewDecoder(w.Body).Decode(&articles)
		validData := len(articles) == 0
		return statusOK && err == nil && validData
	})
}

func TestSearchArticles_Fail(t *testing.T) {
	r := setUpRouter(false)

	// Set up the Gin router to handle the request
	apiGroup := r.Group("/api")
	apiGroup.POST("/search", api.SearchArticles)
	// Create a new HTTP request with form data
	payload := strings.NewReader(`{"keyword":"title"}`)
	req, _ := http.NewRequest("POST", "/api/search", payload)
	req.Header.Set("Content-Type", "application/json")

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusBadRequest
	})
}
