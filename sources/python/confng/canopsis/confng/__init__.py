from pkgutil import extend_path
__path__ = extend_path(__path__, __name__)

from .simpleconf import Configuration, Ini, Json
from .helpers import cfg_to_array