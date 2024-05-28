import os
import sqlite3

import click
from flask import Flask, current_app, g


def get_db():
    if 'db' not in g:
        g.db = sqlite3.connect(
            current_app.config['DATABASE'],
            detect_types=sqlite3.PARSE_DECLTYPES
        )
        g.db.row_factory = sqlite3.Row

    return g.db

def close_db(e=None):
    db = g.pop('db', None)

    if db is not None:
        db.close()

def init_db():
    db = get_db()

    with current_app.open_resource('../../db_init.sql') as f:
        db.executescript(f.read().decode('utf8'))

def delete_db():
    os.remove(current_app.config['DATABASE'])

def query_db(query, args=(), one=False):
    cur = get_db().cursor().execute(query, args)
    rv = cur.fetchall()
    cur.close()
    return (rv[0] if rv else None) if one else rv

def insert_db(query, args=()):
    db = get_db()
    cur = get_db().cursor().execute(query, args)
    db.commit()
    cur.close()
    return cur.lastrowid


@click.command('init-db')
def init_db_command():
    """Clear the existing data and create new tables."""
    init_db()
    click.echo('Initialized the database.')

def init_app(app: Flask):
    app.teardown_appcontext(close_db)
    app.cli.add_command(init_db_command)