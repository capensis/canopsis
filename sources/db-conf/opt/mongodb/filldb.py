#!/usr/bin/env python
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

import subprocess
import sys


if __name__ == '__main__':
    if len(sys.argv) != 2:
        print('Usage: {0} [init|update]'.format(sys.argv[0]))
        sys.exit(1)

    action = sys.argv[1].lower()

    if action not in ['update', 'init']:
        print('Invalid option: {0}'.format(action))
        sys.exit(1)

    subprocess.call('canopsis-filldb --{0}'.format(action), shell=True)
