from flask import Flask, render_template, request, make_response, g
from redis import Redis
import os
import socket
import random
import json
import psycopg2


# ------------------------------
# Helper Functions
# ------------------------------


# Create a connection to the redis server
def create_redis_connection():
    return Redis(host="redis", db=0, socket_timeout=5)

# Create a connection to the postgres server
def create_database_connection():
    return psycopg2.connect("dbname=docker user=docker")


# ------------------------------
# App Setup
# ------------------------------

# Create the flask application
app = Flask(__name__)

# Attach the default "home" route to the flask application
@app.route("/", methods=['GET'])
def hello_world():
    # redis = create_redis_connection()
    hostname = socket.gethostname()

    # Render the empty index.html template
    resp = make_response(render_template(
        'index.html',
        hostname=hostname,
    ))

    # Tag new users with a unique session identifer
    session_id = request.cookies.get('session_id')
    if not session_id:
        session_id = hex(random.getrandbits(64))[2:-1]
        resp.set_cookie('session_id', session_id)

    # Return the generated response
    return resp


# ------------------------------
# Attach the INDEX/SHOW routes for the Symbols table to the flask application
# ------------------------------
@app.route("/symbols", methods=['GET'])
def get_symbols():
    db = create_database_connection()
    cursor = db.cursor()
    cursor.execute("SELECT id, symbol, name from symbols order by symbol desc")
    output = []
    for record in cursor:
        output.append(dict(
            id=record[0],
            symbol=record[1],
            name=record[2]
        ))
    return flask.jsonify(output)

@app.route("/symbols/<id>", methods=['GET'])
def get_symbol(id):
    db = create_database_connection()
    cursor = db.cursor()
    cursor.execute("SELECT id, symbol, name from symbols WHERE id = %s LIMIT 1", (id, ))
    output = cursor.fetchone()
    if not output:
        return flask.jsonify({ message: 'Not found' }), 404
    return flask.jsonify(output)


# ------------------------------
# Boot the flask application & start listening to the given port
# ------------------------------
if __name__ == "__main__":
    port = int(os.getenv('PORT', '80'))
    app.run(
        host='0.0.0.0',
        port=port,
        debug=True,
        threaded=True,
    )
