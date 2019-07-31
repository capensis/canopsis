# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

import re
import os


VERSION_FILE = "VERSION.txt"

CANOPSIS_EDITION_ARG = '--canopsis-edition'
CANOPSIS_STACK_ARG = '--canopsis-stack'
CANOPSIS_VERSION_ARG = '--canopsis-version'

__VERSION_PATTERN = re.compile(r'(\d+\.\d+\.\d+)')


def get_version_file_path(root_path):
    """
    Join root path with version file name.

    :param root_path: `str` root path.
    :return: `str` version file path.
    """
    return os.path.join(root_path, VERSION_FILE)


def parse_version(version_string):
    """
    Parse version string.
    Expects version like ``x.x.x`` where `x` is one o more digits.

    :param version_string: `str` parsing source.
    :return: `str` parsed version or None if couldn't parsed.
    """
    matched = __VERSION_PATTERN.findall(version_string)
    if matched:
        return matched[0]
    else:
        raise ValueError("Can not parse the version inside VERSION.txt.")


def read_version_file(version_file_path):
    """
    Read text file contents and parse Canopsis version from it.

    :param version_file_path: `str` version file path.
    :return: `str` Canopsis version or None if failed.
    """
    with open(version_file_path, 'r') as fp:
        return parse_version(fp.read())
