# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

import os
import sys


class CanopsisUnsupportedEnvironment(Exception):
    pass


def _root_path():
    root = os.environ.get('CPS_PREFIX', sys.prefix)

    if os.environ.get('CPS_BOOTSTRAP', '0') is '1':
        if os.path.isdir(root):
            return root
        else:
            raise CanopsisUnsupportedEnvironment('not a dir: {}'.format(root))

    if os.path.isdir(os.path.join(root, 'etc')):
        pass

    elif os.path.isdir(os.path.join('opt', 'canopsis', 'etc')):
        root = os.path.join(os.path.sep, 'opt', 'canopsis')

    else:
        msg = 'unsupported environment: '
        'cannot safely detect canopsis root path.'
        raise CanopsisUnsupportedEnvironment(msg)

    return root


root_path = _root_path()
