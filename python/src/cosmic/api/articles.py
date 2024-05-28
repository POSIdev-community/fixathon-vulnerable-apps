from flask import render_template, request, session
from ..db import insert_db, query_db
from . import api_bp as bp, views_bp


@bp.route('/articles', methods=['GET'])
def get_articles():
    articles = get_all_articles()
    response = [map_article_view(article) for article in articles]
    return response, 200

@views_bp.route('/article/<articleId>', methods=['GET'])
def get_article(articleId: int):
    article = query_db("""select a.articleId, a.content, a.title, a.userId, u.username as author
                            from Articles a
                            join Users u on u.userId = a.userId
                            where a.articleId = ?""", [articleId], one=True)
    if article is None:
        return {
            "message": "Article not found",
            "data": None,
            "error": "Not Found"
        }, 404
    return render_template('article.html', **map_article_view(article))

@views_bp.route('/article_create', methods=['GET'])
def create_article():
    return render_template('article_create.html')

@bp.route('/article_create', methods=['POST'])
def create_article_post():
    data = request.form
    title = data.get('title')
    content = data.get('content')
    userId = session.get('user_id')
    if not title or not title.strip():
        return {
            "message": "Title is required",
            "data": None,
            "error": "Bad request"
        }, 400
    if not content or not content.strip():
        return {
            "message": "Content is required",
            "data": None,
            "error": "Bad request"
        }, 400
    if not userId:
        return {
            "message": "Unauthorized",
            "data": None,
            "error": "Unauthorized"
        }, 401
    articleId = insert_db("""insert into Articles (title, content, userId) values (?, ?, ?)""", [title, content, userId])
    if articleId == 0 :
        return {
            "message": "Article could not be created",
            "data": None,
            "error": "Internal Server Error"
        }, 500
    return {
        "message": "Article created successfully",
        "data": articleId,
        "error": None
    }, 201




def get_all_articles():
    return query_db("""select a.articleId, a.content, a.title, a.userId, u.username as author
							from Articles a
							join Users u on u.userId = a.userId""")

def map_article_view(article):
    return {
        'articleId': article['articleId'],
        'content': article['content'],
        'title': article['title'],
        'author': article['author'],
        'userId': article['userId'],
        'authorProfileUrl': f'/author/{article["userId"]}'
    }

def get_articles_by_user(userId):
    return query_db("""select a.articleId, a.content, a.title, a.userId, u.username as author
                            from Articles a
                            join Users u on u.userId = a.userId
                            where a.userId = ?""", [userId])