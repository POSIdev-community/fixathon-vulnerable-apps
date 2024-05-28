from flask import Response, current_app, redirect, render_template, request, session
import jwt

from ..db import insert_db, query_db
from . import api_bp as bp
from . import views_bp


def validate_name_and_password(username: str, password: str) -> bool | str:
    if not username or not username.strip():
        return "Username is required"
    if not password or not password.strip():
        return "Password is required"
    return True

@bp.route('/login', methods=['POST'])
def login_post():
    try:
        data = request.form
        username = data.get('username')
        password = data.get('password')
        redirect_to = data.get('redirect_to')
        if not data:
            return {
                "message": "Please provide username and password",
                "data": None,
                "error": "Bad request"
            }, 400
        # validate input
        is_validated = validate_name_and_password(username, password)
        if is_validated is not True:
            return dict(message='Invalid data', data=None, error=is_validated), 400
        user = query_db(f"SELECT * FROM Users WHERE username =  ? and password = ?",
                        [username, password],
                        one=True)
        # Save the user in session
        if user is None:
            return {
                "message": "Invalid username or password",
                "data": None,
                "error": "Unauthorized"
            }, 401
        session.clear()
        session['user_id'] = user['userId']
        session['username'] = user['username']
        #generate jwt token and save it to cookie and session\
        token = jwt.encode({"user_id": user['userId']}, current_app.config["SECRET_KEY"], algorithm="HS256")
        session['token'] = token

        if not redirect_to or not redirect_to.strip():
            redirect_to = '/'
        resp = Response(status=302)
        resp.headers['Location'] = redirect_to
        resp.headers['Authorization'] = f"Bearer {token}"
        resp.set_cookie(current_app.config["JWT_NAME"], token, httponly=True)
        return resp
    except Exception as e:
        return {
            "message": "Something went wrong",
            "data": None,
            "error": str(e)
        }, 500

@views_bp.route('/login', methods=['GET'])
def login_get():
    return render_template('login.html')

@views_bp.route('/register', methods=['GET'])
def register_get():
    return render_template('register.html')

@bp.route('/register', methods=['POST'])
def register_post():
    try:
        data = request.form
        username = data.get('username')
        password = data.get('password')
        if not data:
            return {
                "message": "Please provide username and password",
                "data": None,
                "error": "Bad request"
            }, 400
        # validate input
        is_validated = validate_name_and_password(username, password)
        if is_validated is not True:
            return dict(message='Invalid data', data=None, error=is_validated), 400
        user = query_db("SELECT * FROM Users WHERE username = ?", [username], one=True)
        if user is not None:
            return {
                "message": "User already exists",
                "data": None,
                "error": "Conflict"
            }, 409
        insert_db("INSERT INTO Users (username, password) VALUES (?, ?)", [username, password])
        resp = Response(status=302)
        resp.headers['Location'] = '/login'
        resp.set_cookie('username', username, httponly=True)
        return resp
    except Exception as e:
        return {
            "message": "Something went wrong",
            "data": None,
            "error": str(e)
        }, 500


@bp.route('/logout', methods=['GET'])
def logout():
    session.clear()
    redirect_to = request.args.get('redirect_to')
    if not redirect_to or not redirect_to.strip():
        redirect_to = '/'
    resp = redirect(redirect_to, 302,)
    resp.delete_cookie(current_app.config["JWT_NAME"])
    return resp