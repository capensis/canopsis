from pkgutil import extend_path
__path__ = extend_path(__path__, __name__)

from canopsis.confng.simpleconf import Configuration
from canopsis.confng.vendor import Ini, Json