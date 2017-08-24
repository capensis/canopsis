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

"""engines project building script."""

from canopsis.common.setup import setup

install_requires = [
    'canopsis.common',
    'canopsis.configuration',
    'canopsis.middleware',
    'canopsis.context',
    'canopsis.event',
    'canopsis.schema',
    'canopsis.sla',
    'canopsis.tools',
    'canopsis.monitoring',
    'canopsis.old'
]

setup(
    description='Engines library',
    install_requires=install_requires,
    keywords='engines'
)
