from flask import Flask, render_template, request, make_response, g
from redis import Redis
import os
import socket
import random
import json



# Create a connection to the redis server
def create_redis_connection():
    return Redis(host="redis", db=0, socket_timeout=5)


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


# Boot the flask application & start listening to the given port
if __name__ == "__main__":
    port = int(os.getenv('PORT', '80'))
    app.run(
        host='0.0.0.0',
        port=port,
        debug=True,
        threaded=True,
    )
