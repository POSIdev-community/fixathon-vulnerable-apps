package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"phdays-app/src/api"
	"phdays-app/src/models"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMyProfileTemplate_Success(t *testing.T) {
	r := setUpRouter(true)
	// Set up a mock article and user
	setupArticle()
	setupUser()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs(mockUser.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}).
			AddRow(mockArticle.ArticleId, mockArticle.Content, mockArticle.Title, mockArticle.UserId, mockArticle.Author))

	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.GET("/my_profile", api.MyProfile)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}

	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/my_profile", nil)
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		pageTemplateOK := strings.Index(w.Body.String(), "<title>My Profile</title>") > 0
		articleAuthor := fmt.Sprintf(`<h1>%s</h1>`, mockUser.UserName)
		pageContentOk := strings.Index(w.Body.String(), articleAuthor) > 0

		return statusOK && pageTemplateOK && pageContentOk
	})
}

func TestMyProfileData_Success(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock article and user
	setupArticle()
	setupUser()

	// Set up the expected response
	type responseDataStruct struct {
		UserName string
		Photo    string
		Articles []models.Article
	}
	expectedData := &responseDataStruct{UserName: "testuser", Photo: "../static/profile_photo1.png", Articles: []models.Article{mockArticle}}

	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs(mockUser.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}).
			AddRow(mockArticle.ArticleId, mockArticle.Content, mockArticle.Title, mockArticle.UserId, mockArticle.Author))

	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.GET("/my_profile", api.MyProfile)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/my_profile", nil)

	// Set the accept header to receive a response with json data
	req.Header = map[string][]string{"Accept": {"application/json"}}
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		var resp *responseDataStruct
		err := json.NewDecoder(w.Body).Decode(&resp)
		validData := err == nil && expectedData.UserName == resp.UserName && expectedData.Photo == resp.Photo && len(resp.Articles) == 1 && expectedData.Articles[0] == resp.Articles[0]
		return statusOK && validData
	})
}

func TestMyProfileTemplate_Fail(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock user
	setupUser()

	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WillReturnError(errors.New("error"))
	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.GET("/my_profile", api.MyProfile)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/my_profile", nil)
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		StatusBadRequest := w.Code == http.StatusBadRequest
		validData := strings.Contains(w.Body.String(), "User not found!")
		return StatusBadRequest && validData
	})
}

func TestProfileTemplate_Success(t *testing.T) {
	r := setUpRouter(true)
	// Set up a mock article and user
	setupArticle()
	setupUser()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs(mockUser.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}).
			AddRow(mockArticle.ArticleId, mockArticle.Content, mockArticle.Title, mockArticle.UserId, mockArticle.Author))

	// Set up the Gin router to handle the request
	r.GET("/profile/:userId", api.Profile)
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/profile/1", nil)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		pageTemplateOK := strings.Index(w.Body.String(), "<title>Author Template</title>") > 0
		articleAuthor := fmt.Sprintf(`<h1>%s</h1>`, mockUser.UserName)
		pageContentOk := strings.Index(w.Body.String(), articleAuthor) > 0

		return statusOK && pageTemplateOK && pageContentOk
	})
}

func TestProfileData_Success(t *testing.T) {
	r := setUpRouter(false)
	// Set up a mock article and user
	setupArticle()
	setupUser()

	// Set up the expected response
	type responseDataStruct struct {
		Author   string
		Articles []models.Article
	}
	expectedData := &responseDataStruct{Author: "testuser", Articles: []models.Article{mockArticle}}

	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	mock.ExpectQuery("select a.articleId, a.content, a.title, a.userId, u.username as author").
		WithArgs(mockUser.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"articleId", "content", "title", "userId", "username"}).
			AddRow(mockArticle.ArticleId, mockArticle.Content, mockArticle.Title, mockArticle.UserId, mockArticle.Author))

	// Set up the Gin router to handle the request
	r.GET("/profile/:userId", api.Profile)
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/profile/1", nil)

	// Set the accept header to receive a response with json data
	req.Header = map[string][]string{"Accept": {"application/json"}}

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		var resp *responseDataStruct
		err := json.NewDecoder(w.Body).Decode(&resp)
		validData := err == nil && expectedData.Author == resp.Author && len(resp.Articles) == 1 && expectedData.Articles[0] == resp.Articles[0]
		return statusOK && validData
	})
}

func TestProfileTemplate_Fail(t *testing.T) {
	r := setUpRouter(true)
	// Set up a mock user
	setupUser()

	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnError(errors.New("error"))

	// Set up the Gin router to handle the request
	r.GET("/profile/:userId", api.Profile)
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/profile/1", nil)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		StatusBadRequest := w.Code == http.StatusBadRequest
		validData := strings.Contains(w.Body.String(), "User not found!")
		return StatusBadRequest && validData
	})
}

func TestUploadPhoto_Fail(t *testing.T) {
	r := setUpRouter(false)
	setupUser()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	// create empty multipart writer
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.Close()

	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.POST("/api/profile/upload_photo", api.UploadPhoto)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}
	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/api/profile/upload_photo", buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		StatusBadRequest := w.Code == http.StatusBadRequest
		validData := strings.Contains(w.Body.String(), "http: no such file")
		return StatusBadRequest && validData
	})
}

func TestUploadPhoto_Success(t *testing.T) {
	r := setUpRouter(false)
	setupUser()
	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	// create empty form data
	testfi, _ := os.Stat("./static_test/profile_photo_test.png")
	file, _ := os.Open("./static_test/profile_photo_test.png")
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	assert.NoError(t, mw.WriteField("foo", "bar"))
	w, err := mw.CreateFormFile("profile-photo", "test")
	if assert.NoError(t, err) {
		io.Copy(w, file)
		assert.NoError(t, err)
	}
	mw.Close()

	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.POST("/api/profile/upload_photo", api.UploadPhoto)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}
	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/api/profile/upload_photo", buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusFound := w.Code == http.StatusFound

		fi, err := os.Stat("./src/static/profile_photo1.png")
		e := os.RemoveAll("./src")
		return statusFound && err == nil && fi != nil && e == nil && testfi.Size() > fi.Size()
	})
}

func TestUploadPhotoFromUrl_Fail(t *testing.T) {
	r := setUpRouter(false)
	setupUser()

	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.POST("/api/profile/upload_photo_url", api.UploadPhotoUrl)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}
	// Create a new HTTP request
	url := "/api/static_test/profile_photo_test.png"
	payload := strings.NewReader(fmt.Sprintf("profile-photo-url=%s", url))
	req, _ := http.NewRequest("POST", "/api/profile/upload_photo_url", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		StatusBadRequest := w.Code == http.StatusBadRequest
		return StatusBadRequest
	})
}

func TestUploadPhotoFromUrl_FailSave(t *testing.T) {
	r := setUpRouter(false)
	setupUser()

	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	dir := "src/static"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		assert.NoError(t, err)
	}

	// create a readonly file
	dst, err := os.Create("./src/static/profile_photo1.png")
	dst.Chmod(fs.FileMode(os.O_RDONLY)) // make file readonly
	assert.NoError(t, err)
	dst.Close()

	mockServer := httptest.NewServer(r)
	defer mockServer.Close()

	// set up the gin router to handle the test static file
	r.StaticFile("/api/static_test/profile_photo_test.png", "./static_test/profile_photo_test.png")

	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.POST("/api/profile/upload_photo_url", api.UploadPhotoUrl)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}
	// Create a new HTTP request
	url := fmt.Sprintf("%s/api/static_test/profile_photo_test.png", mockServer.URL)

	payload := strings.NewReader(fmt.Sprintf("profile-photo-url=%s", url))
	req, _ := http.NewRequest("POST", "/api/profile/upload_photo_url", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		StatusBadRequest := w.Code == http.StatusBadRequest
		e := os.RemoveAll("./src")
		return StatusBadRequest && e == nil
	})
}

func TestUploadPhotoFromUrl_FileNotFound(t *testing.T) {
	r := setUpRouter(false)
	setupUser()

	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	mockServer := httptest.NewServer(r)
	defer mockServer.Close()

	// set up the gin router to handle the test static file
	r.StaticFile("/test/static_test/profile_photo_test.png", "./static_test/profile_photo_test.png")

	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.POST("/api/profile/upload_photo_url", api.UploadPhotoUrl)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}
	// Create a new HTTP request
	url := fmt.Sprintf("%s/test/static_test/profile_photo_test1.png", mockServer.URL)
	//"https://uxwing.com/wp-content/themes/uxwing/download/peoples-avatars/corporate-user-icon.png"
	payload := strings.NewReader(fmt.Sprintf("profile-photo-url=%s", url))
	req, _ := http.NewRequest("POST", "/api/profile/upload_photo_url", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		notFound := w.Code == http.StatusNotFound
		return notFound
	})
}

func TestUploadPhotoFromUrl_Success(t *testing.T) {
	r := setUpRouter(false)
	setupUser()

	db, mock, _ := sqlmock.New()
	//models.SetDbConnection() with mocked db
	models.SetDbConnection(db, false)
	defer models.RemoveDbConnection()
	defer db.Close()

	mock.ExpectQuery("SELECT password, username, userId FROM Users").
		WithArgs(fmt.Sprint(mockUser.UserId)).
		WillReturnRows(sqlmock.NewRows([]string{"password", "username", "userId"}).
			AddRow(mockUser.Password, mockUser.UserName, mockUser.UserId))

	mockServer := httptest.NewServer(r)
	defer mockServer.Close()
	// set up the gin router to handle the test static file
	r.StaticFile("/test/static_test/profile_photo_test.png", "./static_test/profile_photo_test.png")

	// Set up the Gin router to handle the request with the authentication middleware
	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.POST("/api/profile/upload_photo_url", api.UploadPhotoUrl)
	}
	tokenString, err := createAuthTokenForTest("1")
	assert.NoError(t, err)
	cookie := &http.Cookie{Name: api.JwtCookie, Value: tokenString}
	// Create a new HTTP request
	url := fmt.Sprintf("%s/test/static_test/profile_photo_test.png", mockServer.URL)

	payload := strings.NewReader(fmt.Sprintf("profile-photo-url=%s", url))
	req, _ := http.NewRequest("POST", "/api/profile/upload_photo_url", payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	checkHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusFound := w.Code == http.StatusFound

		fi, err := os.Stat("./src/static/profile_photo1.png")
		e := os.RemoveAll("./src")
		validRedirect := w.Header().Get("Location") == "/my_profile"
		return statusFound && err == nil && fi != nil && e == nil && validRedirect
	})
}
