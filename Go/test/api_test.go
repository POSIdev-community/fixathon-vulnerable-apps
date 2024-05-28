package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"phdays-app/src/api"
	"phdays-app/src/models"
)

var mockUser models.User
var mockArticle models.Article

func setupUser() {
	mockUser = models.User{
		UserId:   1,
		UserName: "testuser",
		Password: "testpassword",
	}
}

func setupArticle() {
	mockArticle = models.Article{
		Title:            "title",
		Content:          "content",
		ArticleId:        1,
		UserId:           1,
		Author:           "testuser",
		AuthorProfileUrl: "/profile/1",
	}
}

func TestLogin_Success(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Set up a mock user
	setupUser()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()

	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(mockUser.UserName, mockUser.Password).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	// Create a new HTTP request with form data
	payload := strings.NewReader("username=testuser&password=testpassword&redirect_to=/dashboard")
	req, _ := http.NewRequest("POST", "/login", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Set up the Gin router to handle the request
	router.POST("/login", api.Login)

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert that the response status code is 302
	assert.Equal(t, http.StatusFound, w.Code)

	// Assert that the JWT cookie is set
	assert.Contains(t, w.Header().Get("Set-Cookie"), api.JwtCookie)

	// Assert that the redirect URL is correct
	assert.Equal(t, "/dashboard", w.Header().Get("Location"))
}

func TestLogin_InvalidParameters(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Create a new HTTP request with empty form data
	payload := strings.NewReader("username=&password=&redirect_to=")
	req, _ := http.NewRequest("POST", "/login", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Set up the Gin router to handle the request
	router.POST("/login", api.Login)

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert that the response status code is 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Assert that the error message is correct
	assert.Contains(t, w.Body.String(), "Parameters can't be empty")
}

func TestLogin_AuthenticationFailed(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Set up a mock user
	setupUser()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()

	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(mockUser.UserName, mockUser.Password).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}))

	// Create a new HTTP request with valid form data
	payload := strings.NewReader("username=testuser&password=testpassword&redirect_to=/dashboard")
	req, _ := http.NewRequest("POST", "/login", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Set up the Gin router to handle the request
	router.POST("/login", api.Login)

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert that the response status code is 401 (Unauthorized)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Assert that the error message is correct
	assert.Contains(t, w.Body.String(), "Authentication failed")
}

func TestLogout_Success(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Create a new HTTP request with form data
	req, _ := http.NewRequest("POST", "/logout", nil)

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Set up the Gin router to handle the request
	router.POST("/logout", api.Logout)

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert that the response status code is 302
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the JWT cookie is set
	assert.Contains(t, w.Header().Get("Set-Cookie"), api.JwtCookie)

	// Assert that the JWT cookie is empty
	assert.Equal(t, w.Result().Cookies()[0].Name, api.JwtCookie)
	assert.Equal(t, w.Result().Cookies()[0].Value, "")

	// Assert that the result message is correct
	assert.Contains(t, w.Body.String(), "You are logged out")
}

func TestIndex_Success(t *testing.T) {
	r := setUpRouter(true)
	r.GET("/", api.Index)
	req, _ := http.NewRequest("GET", "/", nil)
	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		pageOK := strings.Index(w.Body.String(), "<h1>Cosmic Blogs</h1>") > 0

		return statusOK && pageOK
	})
}

func TestLoginPage_Success(t *testing.T) {
	r := setUpRouter(true)
	r.GET("/login", api.LoginPage)
	req, _ := http.NewRequest("GET", "/login", nil)
	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		pageOK := strings.Index(w.Body.String(), "<h1>Welcome to Cosmic Blogs</h1>") > 0

		return statusOK && pageOK
	})
}

func TestAuth_ArticleCreateTemplate_Success(t *testing.T) {
	r := setUpRouter(true)
	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.GET("/article_create", api.ArticleCreate)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}

	req, _ := http.NewRequest("GET", "/article_create", nil)
	req.AddCookie(cookie)

	// Assert that the response contains article_template
	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		pageOK := strings.Index(w.Body.String(), "<title>Create article form</title>") > 0

		return statusOK && pageOK
	})
}

func TestAuth_ArticleCreateTemplate_NoToken(t *testing.T) {
	r := setUpRouter(true)
	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.GET("/article_create", api.ArticleCreate)
	}

	req, _ := http.NewRequest("GET", "/article_create", nil)

	// Assert that the response is redirected to login page
	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusFound := w.Code == http.StatusFound
		validRedirect := w.Header().Get("Location") == "/login"
		return statusFound && validRedirect
	})
}

func TestAuth_ArticleCreateTemplate_InvalidToken(t *testing.T) {
	r := setUpRouter(true)
	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.GET("/article_create", api.ArticleCreate)
	}

	// create the invalid token
	tokenString := "test"
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}

	req, _ := http.NewRequest("GET", "/article_create", nil)
	req.AddCookie(cookie)

	// Assert that the response is redirected to login page
	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusFound := w.Code == http.StatusFound
		validRedirect := w.Header().Get("Location") == "/login"
		return statusFound && validRedirect
	})
}

// Helper function to process a request and test its response
func checkHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func setUpRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("../src/html/*")
	}
	return r
}

// Function to create JWT tokens with claims for tests
func createAuthTokenForTest(userId string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,                           // Subject (user identifier)
		"iss": "phdays-app",                     // Issuer
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	var secret = []byte("secret")
	tokenString, err := claims.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
