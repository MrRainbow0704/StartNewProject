from flask import render_template, Response
from . import app


@app.route("/")
def home() -> Response:
    return render_template("home.html"), 200
