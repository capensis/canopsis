# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

"""
Here are all url converters used to convert/validate bottle url parameters.
"""

from __future__ import unicode_literals

from bottle import abort
import json
import logging

logger = logging.getLogger('webserver')


def mongo_filter(config):
    """
    Parse mongo filter format from url
    """
    regexp = r'{.+}'

    def to_python(match):
        """
        Convert and validate the url parameter to python
        """
        try:
            return json.loads(match)
        except:
            logger.error('Cannot parse url parameter: {}'.format(match))
            return abort(400)

    def to_url(filter_):
        """
        Convert json object to url format
        """
        return filter_  # a simple dict !

    return regexp, to_python, to_url


def id_filter(config):
    """
    Parse an id from url
    """
    regexp = r'.+'

    def to_python(match):
        """
        Convert and validate the url parameter to python
        """
        # TODO: do more stuff, like searching the ID into db
        return match

    def to_url(filter_):
        """
        Convert json object to url format
        """
        return filter_

    return regexp, to_python, to_url
