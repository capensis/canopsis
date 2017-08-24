import os
import sys


class CanopsisUnsupportedEnvironment(Exception):
    pass

def _root_path():
    root = None

    if os.path.isdir(os.path.join(sys.prefix, 'etc')):
        root = sys.prefix

    elif os.path.isdir(os.path.join('opt', 'canopsis', 'etc')):
        root = os.path.join(os.path.sep, 'opt', 'canopsis')

    else:
        msg = 'unsupported environment: cannot safely detect canopsis root path.'
        raise CanopsisUnsupportedEnvironment(msg)

    return root

root_path = _root_path()
