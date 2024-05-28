package api

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"phdays-app/src/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var profilePhoto = "/static/profile_photo"

func MyProfile(c *gin.Context) {
	id := c.GetString(loggedInUser)
	user, err := models.GetUserById(id)
	if err != nil {
		logrus.WithContext(context.Background()).Errorf("User not found %v", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	articles := models.GetArticlesByUserId(user.UserId)
	render(c, "my_profile.html", gin.H{
		"userName": user.UserName,
		"photo":    fmt.Sprintf("..%s%v.png", profilePhoto, user.UserId),
		"articles": articles,
	})
}

func Profile(c *gin.Context) {
	id := c.Param("userId")

	author, err := models.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	articles := models.GetArticlesByUserId(author.UserId)
	render(c, "author_template.html", gin.H{
		"author":   author.UserName,
		"articles": articles,
	})
}

func UploadPhoto(c *gin.Context) {
	file, err := c.FormFile("profile-photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetString(loggedInUser)
	fileDst := fmt.Sprintf("./src%s%v.png", profilePhoto, userId)
	// Upload the file to specific dst.
	err = c.SaveUploadedFile(file, fileDst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	path, _ := filepath.Abs(fileDst)
	runImageMagicCmd(path)

	c.Redirect(http.StatusFound, "/my_profile")
}

func UploadPhotoUrl(c *gin.Context) {

	url := c.PostForm("profile-photo-url")
	//https://uxwing.com/wp-content/themes/uxwing/download/peoples-avatars/corporate-user-icon.png

	response, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if response.StatusCode == http.StatusNotFound {
		c.JSON(response.StatusCode, gin.H{"error": response.Status})
		return
	}
	defer response.Body.Close()

	err = saveFile(response.Body, c.GetString(loggedInUser))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/my_profile")
}

func runImageMagicCmd(command string) {
	var cmd = &exec.Cmd{}
	//cmd = exec.Command("magick", "mogrify", "-format", "jpg", "-quality", "50", command)
	cmd = exec.Command("magick", "mogrify", "-resize", "50%%", command)
	err := cmd.Run()
	if err != nil {
		logrus.Info(err.Error())
	}
	logrus.Info(cmd.String())
}

func saveFile(source io.ReadCloser, userId string) error {

	dstFile := fmt.Sprintf("./src%s%v.png", profilePhoto, userId)
	dir := filepath.Dir(dstFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//open a file for writing
	destination, err := os.OpenFile(dstFile, os.O_WRONLY, 0644)
	if err != nil {
		destination, err = os.Create(fmt.Sprintf("./src%s%v.png", profilePhoto, userId))
	}

	if err != nil {
		return err
	}
	defer destination.Close()
	src, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}
	_, err = destination.Write(src)
	if err != nil {
		return err
	}

	return nil
}
