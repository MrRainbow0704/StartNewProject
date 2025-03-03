from flask import Flask
import config
from . import tools


app = Flask(__name__, root_path=config.ROOT_DIR)
app.secret_key = config.SECRET_KEY

from . import routes
