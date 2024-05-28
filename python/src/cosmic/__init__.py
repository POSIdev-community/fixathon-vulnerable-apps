import os

from flask import Flask, render_template, session
from .api import views_bp, api_bp
from . import config


def create_app(test_config=None):
    # create and configure the app
    app = Flask(__name__, instance_relative_config=True, static_url_path='/static', static_folder='../static')
    app.config.from_mapping(
        config.config,
        DATABASE=os.path.join(app.instance_path, 'cosmic_db.sqlite'),
    )

    if not os.path.isdir("static"):
        os.mkdir("static")
    if not os.path.isdir("temp"):
        os.mkdir("temp")        

    if test_config is None:
        # load the instance config, if it exists, when not testing
        app.config.from_pyfile('config.py', silent=True)
    else:
        # load the test config if passed in
        app.config.from_mapping(test_config)

    # ensure the instance folder exists
    try:
        os.makedirs(app.instance_path)
    except OSError:
        pass

    # a simple page that says hello
    app.register_blueprint(views_bp)
    app.register_blueprint(api_bp)
    
    from . import db
    db.init_app(app)

    return app

if __name__ == "__main__":
    app = create_app()
    app.run(port=8080, debug=True)