import os
from flask import Blueprint, Response, redirect, render_template, request, session

from ..magick import convert_to_png
from .articles import get_articles_by_user
from .token import token_required
from ..db import query_db
from . import api_bp as bp, views_bp
import requests
from PIL import Image
from io import BytesIO


@views_bp.route('/my_profile', methods=['GET'])
@token_required
def my_profile():
    user = query_db("SELECT * FROM Users WHERE userId = ?", [session['user_id']], one=True)
    if user is None:
        return {
            "message": "User not found",
            "data": None,
            "error": "Not Found"
        }, 404
    user = {
        'userId': user['userId'],
        'username': user['username'],
        'articles': get_articles_by_user(user['userId'])
    }
    return render_template('my_profile.html', **user)


@bp.route('/profile/upload_photo_url', methods=['POST'])
@token_required
def upload_photo_url():    
    """Uploads a photo for user by given url and save it to static/profile_photo{userId}.png"""
    data = request.form
    photo_url = data.get('profile-photo-url')
    if not photo_url:
        return {
            "message": "Photo url is required",
            "data": None,
            "error": "Bad request"
        }, 400
    user = query_db("SELECT * FROM Users WHERE userId = ?", [session['user_id']], one=True)
    if user is None:
        return {
            "message": "User not found",
            "data": None,
            "error": "Not Found"
        }, 404

    #download photo by photo_url
    response = requests.get(photo_url)
    try:
        img = Image.open(BytesIO(response.content))
        img.save(f"static/profile_photo{user['userId']}.png")
        return redirect('/my_profile', 302)
    except Exception as e:
        return {
            "message": "Photo could not be uploaded",
            "data": None,
            "error": "Internal Server Error"
        }, 500

@bp.route('/profile/upload_photo', methods=['POST'])
def upload_photo():
    """Uploads a photo for user and save it to static/profile_photo{userId}.png"""
    data = request.files
    photo = data.get('profile-photo')
    if not photo:
        return {
            "message": "Photo is required",
            "data": None,
            "error": "Bad request"
        }, 400
    user = query_db("SELECT * FROM Users WHERE userId = ?", [session['user_id']], one=True)
    if user is None:
        return {
            "message": "User not found",
            "data": None,
            "error": "Not Found"
        }, 404
    
    # save photo to temp dir and then convert it to png from jpg with magick mogrify
    temp_file_name = f"temp/{photo.filename}"
    if os.path.exists(temp_file_name):
        os.remove(temp_file_name)
    f = open(temp_file_name, 'wb')
    f.write(photo.read())
    f.close()
    p = convert_to_png(temp_file_name)
    
    if p.returncode != 0:
        return {
            "message": "Photo could not be uploaded",
            "data": None,
            "error": "Internal Server Error"
        }, 500
    # Move converted photo to final destination
    photopath = f"static/profile_photo{user['userId']}.png"
    if os.path.exists(photopath):
        os.remove(photopath)
    os.rename(temp_file_name, photopath)
    # upload photo back from temp and save it to static/profile_photo{userId}.png
    return redirect('/my_profile', 302)

