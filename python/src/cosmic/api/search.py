

from flask import Response, render_template, request

from .articles import map_article_view
from ..db import query_db
from . import api_bp, views_bp


@views_bp.route('/search', methods=['GET'])
def search():
    return render_template('search.html')

@api_bp.route('/search', methods=['POST'])
def search_post():
    data = request.form
    query = data.get('query')
    if not query or not query.strip():
        return {
            "message": "Query is required",
            "data": None,
            "error": "Bad request"
        }, 400
    articles = query_db(f"""SELECT a.articleId, a.content, a.title, a.userId, u.username as author
                         from Articles a join Users u on u.userId = a.userId 
                         WHERE title LIKE '{query}' OR content LIKE '{query}'""", one=False)

    response = [map_article_view(article) for article in articles]
    return response, 200