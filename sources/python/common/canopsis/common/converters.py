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
All url converters used to convert/validate bottle url parameters.
"""

from __future__ import unicode_literals

import json
import logging

logger = logging.getLogger('webserver')


def mongo_filter(config):
    """
    Parse mongo filter format from url.
    """
    regexp = r'{.*}'

    def to_python(match):
        """
        Convert and validate the url parameter to python

        :param match: the matched portion of the url with regexp
        :type match: str
        :rtype: dict
        """
        try:
            return json.loads(match)
        except:
            logger.error('Cannot parse url parameter: {}'.format(match))
            raise

    def to_url(filter_):
        """
        Convert json object to url format

        :param filter_: the object to convert to an url
        :type filter_: dict
        :rtype: str
        """
        return filter_  # a simple dict !

    return regexp, to_python, to_url


def id_filter(config):
    """
    Parse a generic id.
    """
    regexp = r'.+'

    def to_python(match):
        """
        Convert and validate the url parameter to python.

        :rtype: str
        """
        # TODO: do more stuff, like searching the ID into db
        return match

    def to_url(filter_):
        """
        Convert json object to url format.
        """
        return filter_

    return regexp, to_python, to_url
