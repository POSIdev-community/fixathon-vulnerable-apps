from functools import wraps
import jwt
from flask import request, abort, session
from flask import current_app

def get_jwt_token() -> str | None:
    jwt_name = current_app.config["JWT_NAME"]
    if jwt_name in request.cookies:
        return request.cookies[jwt_name]
    return None



def token_required(f):
    @wraps(f)
    def decorated(*args, **kwargs):
        token = get_jwt_token()        
        if not token:
            return {
                "message": "Authentication Token is missing!",
                "data": None,
                "error": "Unauthorized"
            }, 401
        try:
            data=jwt.decode(token, current_app.config["SECRET_KEY"], algorithms=["HS256"])
            session_id = session.get('user_id')
            token_id = data.get('user_id')
            if not session_id or not token_id or session_id != token_id:
                return {
                    "message": "Invalid Authentication token!",
                    "data": None,
                    "error": "Unauthorized"
                }, 401
        except Exception as e:
            return {
                "message": "Something went wrong",
                "data": None,
                "error": str(e)
            }, 500

        return f(*args, **kwargs)

    return decorated