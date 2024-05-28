import unittest
from unittest.mock import patch
from . import BaseCosmicTest

#test for /author/<authorId>
class TestAuthor(BaseCosmicTest):

    def test_get_author(self):
        author_id = 1  # replace with an actual author ID from your database
        response = self.client.get(f'/author/{author_id}')
        self.assertEqual(response.status_code, 200)

    def test_get_author_not_found(self):
        author_id = 999  # replace with an ID that doesn't exist in your database
        response = self.client.get(f'/author/{author_id}')
        self.assertEqual(response.status_code, 404)

if __name__ == '__main__':
    unittest.main()