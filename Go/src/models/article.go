package models

import (
	"database/sql"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/sirupsen/logrus"
)

const profileUrl = "/profile"

type Article struct {
	Title            string `json:"title"`
	Content          string `json:"content"`
	ArticleId        int    `json:"id"`
	UserId           int    `json:"userId"`
	Author           string `json:"author"`
	AuthorProfileUrl string `json:"authorProfileUrl"`
}

func GetAllArticles() []Article {

	initDbConnection()
	defer closeDbConnection()
	rows, err := db.Query(`select a.articleId, a.content, a.title, a.userId, u.username as author
							from Articles a
							join Users u on u.userId = a.userId`)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	articles := mapRows(rows)
	return articles
}

func GetArticle(id string) (*Article, error) {

	initDbConnection()
	defer closeDbConnection()
	article := &Article{}
	results, err := db.Query(`select a.articleId, a.content, a.title, a.userId, u.username as author
								from Articles a
								join Users u on u.userId = a.userId where articleId = ?`, id)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer results.Close()

	if results.Next() {
		err = results.Scan(&article.ArticleId, &article.Content, &article.Title, &article.UserId, &article.Author)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return article, err
}

func GetArticlesByUserId(userId int) []Article {

	initDbConnection()
	defer closeDbConnection()
	rows, err := db.Query(`select a.articleId, a.content, a.title, a.userId, u.username as author
							from Articles a
							join Users u on u.userId = a.userId 
							where u.userId = ?`, userId)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	articles := mapRows(rows)
	return articles
}

func GetArticlesByKeyword(keyword string) []Article {
	searchString := "%" + keyword + "%"
	initDbConnection()
	defer closeDbConnection()
	query := `select a.articleId, a.content, a.title, a.userId, u.username as author
				from Articles a
				join Users u on u.userId = a.userId 
				where title like ? or content like ?`

	rows, err := db.Query(query, searchString, searchString)
	if err != nil {
		logrus.Error(err)
	}
	defer rows.Close()
	articles := mapRows(rows)
	return articles
}

func AddNewArticle(content, title, userId string) (int64, error) {
	initDbConnection()
	defer closeDbConnection()
	query := fmt.Sprintf(`INSERT INTO articles (content, title, userId) VALUES (%s,%s,%s)`, content, title, userId)
	insert, err := db.Exec(query)

	// if there is an error inserting, handle it
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return insert.LastInsertId()
}

func mapRows(rows *sql.Rows) []Article {

	articles := []Article{}
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ArticleId, &article.Content, &article.Title, &article.UserId, &article.Author)
		if err != nil {
			logrus.Error(err)
		}

		article.AuthorProfileUrl = fmt.Sprintf("%s/%v", profileUrl, article.UserId)
		articles = append(articles, article)
	}
	err := rows.Err()
	if err != nil {
		logrus.Error(err)
	}

	return articles
}
