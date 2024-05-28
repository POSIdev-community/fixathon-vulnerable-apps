import os
from pathlib import Path
import unittest
from unittest.mock import patch

import werkzeug
from . import BaseCosmicTest
#test for /my_profile, api/profile/upload_photo_url, api/profile/upload_photo
class TestProfile(BaseCosmicTest):

    def test_get_my_profile_unauthorized(self):
        response = self.client.get('my_profile')
        self.assertEqual(response.status_code, 401)

    def test_get_my_profile(self):
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/login', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.get('my_profile')
        self.assertEqual(response.status_code, 200)

    def test_upload_photo_url_invalid(self):
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/login', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/profile/upload_photo_url', data={'profile-photo-url': 'https://example.com/photo.jpg'})
        self.assertEqual(response.status_code, 500)
    
    def test_upload_photo_url(self):
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/login', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/profile/upload_photo_url', data={'profile-photo-url': 'https://raw.githubusercontent.com/mathiasbynens/small/master/jpeg.jpg'})
        self.assertEqual(response.status_code, 302)


    def test_upload_photo(self):
        response = self.client.post('api/register', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        response = self.client.post('api/login', data=dict(username='test', password='test'))
        self.assertEqual(response.status_code, 302)
        f = open(f'{Path(__file__).parent}/example_img.png', "rb")
        file = werkzeug.datastructures.FileStorage(
            stream=f,
            filename="example_img.png",
            content_type="image/png",
        )
        response = self.client.post(
            'api/profile/upload_photo',
            data={
                'profile-photo': file,
            },
            follow_redirects=True,
            content_type='multipart/form-data',
        )
        f.close()
        if not os.path.exists(f"static/profile_photo4.png"):
            self.fail("File not saved")
        os.remove(f"static/profile_photo4.png")
        
        self.assertEqual(response.status_code, 200)

if __name__ == '__main__':
    unittest.main()