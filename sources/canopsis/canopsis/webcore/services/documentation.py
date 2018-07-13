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

from __future__ import unicode_literals

import codecs
import markdown
import os

from canopsis.common import root_path

DOC_DIR = os.path.join(root_path, "var/www/src/documentations")

PAGE_TEMPLATE = "<html><title>Documentation</title><body>{0}</body></html>"
IO_ERROR = "<p>Can not render the page : documentation not found</p>"
OTHER_ERROR = "<p>Can not render the documentation requested.</p>"


def get_HTML_content(name):
    """Generate and return as HTML page the documentation requested

    Use this method carefully, the markdown package does not list the
    exceptions that can be thrown.

    :param name: the documentation page name
    :return : the documentation as HTML
    """
    name = name + ".md"
    filename = os.path.join(DOC_DIR, name)

    with codecs.open(filename, mode="r", encoding="utf-8") as fd:
        lines = fd.readlines()

    content = ""
    for line in lines:
        content += line

    html = markdown.markdown(content)

    return PAGE_TEMPLATE.format(html)


def exports(ws):

    @ws.application.route(
        "/api/v2/documentation/<name>"
    )
    def get_documentation_page(name):
        """
        Return a documentation page.
        :param the name of the documentation
        :return: the documentation as HTML
        """

        try:
            return get_HTML_content(name)
        except IOError:
            ws.logger.error("Documentation page {} not found.".format(name))
            return PAGE_TEMPLATE.format(IO_ERROR)
        except Exception:
            msg = "An error occured during the generation of the documentation."
            ws.logger.exception(msg)
            return PAGE_TEMPLATE.format(OTHER_ERROR)
