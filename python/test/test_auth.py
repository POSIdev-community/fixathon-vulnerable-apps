import unittest
from unittest.mock import patch
from . import BaseCosmicTest

#test for api/login, api/register, api/logout
class TestAuth(BaseCosmicTest):

    def test_register(self):
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)

    def test_register_no_data(self):
        response = self.client.post('api/register')
        self.assertEqual(response.status_code, 400)

    def test_register_user_exists(self):
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 409)

    def test_login(self):
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/login', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)

    def test_login_no_data(self):
        response = self.client.post('api/login')
        self.assertEqual(response.status_code, 400)

    def test_login_user_not_found(self):
        response = self.client.post('api/login', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 401)

    def test_logout(self):
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/login', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.get('api/logout')
        self.assertEqual(response.status_code, 302)

if __name__ == '__main__':
    unittest.main()