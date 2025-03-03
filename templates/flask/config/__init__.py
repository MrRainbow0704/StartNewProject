from pathlib import Path


ROOT_DIR = Path(__file__).resolve().parent.parent
HOST = "0.0.0.0"
PORT = 8080
DEBUG = True
SECRET_KEY = "ThisShouldBeSecret"

from .gunicorn import *
