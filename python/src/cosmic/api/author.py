from flask import render_template
from .articles import get_articles_by_user
from ..db import query_db
from . import api_bp as bp, views_bp

@views_bp.route('/author/<authorId>', methods=['GET'])
def get_author(authorId: int):
    user = query_db("SELECT * FROM Users WHERE userId = ?", [authorId], one=True)
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
    return render_template('author.html', **user)