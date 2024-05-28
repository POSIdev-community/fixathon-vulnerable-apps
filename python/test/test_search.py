import unittest
from unittest.mock import patch
from . import BaseCosmicTest

#test for api/search, /search
class TestSearch(BaseCosmicTest):    
        def test_search(self):
            response = self.client.get('search')
            self.assertEqual(response.status_code, 200)
    
        def test_search_post_no_data(self):
            response = self.client.post('api/search')
            self.assertEqual(response.status_code, 400)
    
        def test_search_post(self):
            response = self.client.post('api/search', data=dict(query='test'))
            self.assertEqual(response.status_code, 200)