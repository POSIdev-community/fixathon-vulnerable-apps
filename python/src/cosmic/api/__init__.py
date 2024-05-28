from flask import (
    Blueprint, render_template, session
)

api_bp = Blueprint('auth', __name__, url_prefix='/api')
views_bp = Blueprint('views', __name__)

#init all views
from . import articles, auth, author, profile, search

@views_bp.route('/')
def index():
    return render_template('index.html', **{'authorized': session.get('user_id') is not None})

# @views_bp.route('/static/<path>', methods=['GET'])
# def send_report(path):
#     return send_from_directory('../static', path)