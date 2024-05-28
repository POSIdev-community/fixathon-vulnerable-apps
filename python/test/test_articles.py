
import unittest
from unittest.mock import patch
from . import BaseCosmicTest

class TestGetArticles(BaseCosmicTest):

    def test_get_articles(self):
        response = self.client.get('api/articles')
        firstArticle = response.json[0]
        self.assertEqual(response.status_code, 200)
        self.assertEqual(firstArticle['articleId'], 1)
        self.assertEqual(firstArticle['author'], 'GalacticExplorer')
        self.assertEqual(firstArticle['authorProfileUrl'], '/author/1')
        self.assertEqual(firstArticle['content'], 'Наше путешествие в Галактику Андромеды раскрыло захватывающие виды и таинственные явления.')
        self.assertEqual(firstArticle['title'], 'Исследование Галактики Андромеды')
        self.assertEqual(firstArticle['userId'], 1)

    def test_get_article(self):
        article_id = 1  # replace with an actual article ID from your database
        response = self.client.get(f'/article/{article_id}')
        self.assertEqual(response.status_code, 200)


    def test_get_article_not_found(self):
        article_id = 999  # replace with an ID that doesn't exist in your database
        response = self.client.get(f'/article/{article_id}')
        self.assertEqual(response.status_code, 404)

if __name__ == '__main__':
    unittest.main()