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
from canopsis.engines.core import TaskHandler
from canopsis.old.account import Account
from canopsis.old.storage import Storage
from canopsis.old.file import File
from canopsis.common.utils import ensure_unicode
from canopsis.common.template import Template

class engine(TaskHandler):
    etype = 'taskfile'

    def handle_task(self, job):
        path = job.get('path', '~/tmp/{{timestamp}}_{{rk}}')
        body = job.get('body', None)

        template_data = job.get('jobctx', {})
        path = Template(path)(template_data)
        path = os.path.expanduser(path)
        body = Template(body)(template_data)

        try:
            with open(path, 'a+') as f:
                f.write(body)
            f.close
            return (0, "File has been write successfully")
        except Exception as err:
            self.logger.error("Can not write file in {0}: {1}".format(path, err))
            return (2, "Can not write file in {0}: {1}".format(path, err))
