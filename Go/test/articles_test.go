package api_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"phdays-app/src/api"
	"phdays-app/src/models"
)

func TestArticleTemplate_Success(t *testing.T) {
	r := setUpRouter(true)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs(fmt.Sprint(mockArticle.ArticleId)).
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}).
			AddRow(mockArticle.ArticleId, mockArticle.Content, mockArticle.Title, mockArticle.UserId, mockArticle.Author))

	// Set up the Gin router to handle the request
	r.GET("/articles/:id", api.Article)
	req, _ := http.NewRequest("GET", "/articles/1", nil)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		pageArticleTemplateOK := strings.Index(w.Body.String(), "<title>Article Template</title>") > 0
		articleAuthor := fmt.Sprintf(`<p><a href="../profile/%v">%s</a></p>`, mockArticle.UserId, mockArticle.Author)
		pageContentOk := strings.Index(w.Body.String(), articleAuthor) > 0

		return statusOK && pageArticleTemplateOK && pageContentOk
	})
}

func TestArticleTemplate_Fail(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs(fmt.Sprint(mockArticle.ArticleId)).
		WillReturnError(errors.New("error"))

	// Set up the Gin router to handle the request
	r.GET("/articles/:id", api.Article)
	req, _ := http.NewRequest("GET", "/articles/1", nil)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNotFound
	})
}

func TestArticleData_Success(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs(fmt.Sprint(mockArticle.ArticleId)).
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}).
			AddRow(mockArticle.ArticleId, mockArticle.Content, mockArticle.Title, mockArticle.UserId, mockArticle.Author))

	// Set up the Gin router to handle the request
	r.GET("/articles/:id", api.Article)
	req, _ := http.NewRequest("GET", "/articles/1", nil)
	// Set the accept header to receive a response with json data
	req.Header = map[string][]string{"Accept": {"application/json"}}

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		var article models.Article
		err := json.NewDecoder(w.Body).Decode(&article)
		validData := article.Author == mockArticle.Author && article.Content == mockArticle.Content && article.Title == mockArticle.Title && article.UserId == mockArticle.UserId
		return statusOK && err == nil && validData
	})
}

func TestApiArticles_Success(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}).
			AddRow(mockArticle.ArticleId, mockArticle.Content, mockArticle.Title, mockArticle.UserId, mockArticle.Author))

	// Set up the Gin router to handle the request
	apiGroup := r.Group("/api")
	apiGroup.GET("/articles", api.Articles)
	req, _ := http.NewRequest("GET", "/api/articles", nil)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		var articles []models.Article
		err := json.NewDecoder(w.Body).Decode(&articles)
		validData := articles[0] == mockArticle
		return statusOK && err == nil && validData
	})
}

func TestApiArticles_Fail(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}))

	// Set up the Gin router to handle the request
	r.GET("api/articles", api.Articles)
	req, _ := http.NewRequest("GET", "api/articles", nil)

	// Assert that the response status code is 404
	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusNotFound := w.Code == http.StatusNotFound

		return statusNotFound
	})
}

func TestArticleCreateTemplate_Success(t *testing.T) {
	r := setUpRouter(true)
	r.GET("/article_create", api.ArticleCreate)
	req, _ := http.NewRequest("GET", "/article_create", nil)

	// Assert that the response contains article_template
	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		pageOK := strings.Index(w.Body.String(), "<title>Create article form</title>") > 0

		return statusOK && pageOK
	})
}

func TestCreateArticle_Success(t *testing.T) {
	// Create a new Gin router
	router := setUpRouter(false)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()

	defer db.Close()

	mock.ExpectExec("INSERT INTO articles").
		WillReturnResult(sqlmock.NewResult(int64(mockArticle.ArticleId), 1))

	// Set up the Gin router to handle the request with the authentication middleware
	gr := router.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.POST("/api/article_create", api.CreateArticle)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}

	// Create a new HTTP request with form data
	payload := strings.NewReader("title=title&content=content&redirect_to=/dashboard")
	req, _ := http.NewRequest("POST", "/api/article_create", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	checkHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		// Assert that the response status code is 302
		statusFound := w.Code == http.StatusFound

		// Assert that the redirect URL is correct
		validRedirect := fmt.Sprintf("/articles/%v", mockArticle.ArticleId) == w.Header().Get("Location")

		return statusFound && validRedirect
	})
}

func TestCreateArticle_Fail(t *testing.T) {
	// Create a new Gin router
	router := setUpRouter(false)
	// Set up a mock article
	setupArticle()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()

	defer db.Close()

	mock.ExpectExec("INSERT INTO articles").
		WillReturnError(fmt.Errorf("error"))

		// Set up the Gin router to handle the request with the authentication middleware
	gr := router.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.POST("/api/article_create", api.CreateArticle)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}

	// Create a new HTTP request with form data
	payload := strings.NewReader("title=title&content=content&redirect_to=/dashboard")
	req, _ := http.NewRequest("POST", "/api/article_create", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	checkHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		// Assert that the response status code is 500
		return w.Code == http.StatusInternalServerError
	})
}
