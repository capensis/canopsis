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

"""
Task manager module.

Permits to load and save dynamically tasks in a distributed environment.
"""

from sys import path

from shutil import copy

from importlib import import_module
try:  # PYTHON3
    from importlib import reload
except ImportError:
    pass  # PYTHON2

from canopsis.common.utils import lookup
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths)

CONF_PATH = 'task/task.conf'
CATEGORY = 'TASK'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class TaskManager(MiddlewareRegistry):
    """
    TaskManager manage dedicated to python tasks.

    Contains a local dictionary of {tasks name: (task, task_info)} which is
    updated at runtime depending of the distributed environment and other tasks
    executed on other machines::

        A task is a python callable object, and a task_info is a dictionary
        which could be gave to the task object such as a kwargs.

    Interacts with both storage:

    - default storage: contains task information.
    - file storage: contains additional task source files.

    In order to load dynamically tasks from the file storage, it also manage a
    task directory which contains all task source files.
    """

    NAME = 'task'  #: task field name

    STORAGE = 'task_storage'  #: storage item name
    FILE_STORAGE = 'task_file_storage'  #: file storage item name

    def __init__(self, tasks=None, task_directory=None, *args, **kwargs):

        super(TaskManager, self).__init__(*args, **kwargs)

        self.tasks = {} if tasks is None else tasks
        self.task_directory = task_directory

    @property
    def task_directory(self):
        """
        :return: self task directory which contains task files.
        :rtype: str
        """

        return self._task_directory

    @task_directory.setter
    def task_directory(self, value):
        """
        Change of task directory.

        :param str value: new task directory path to use.
        """

        if value is not None:
            # add value in path if not None
            if value not in path:
                path.append(value)
        self._task_directory = value

    @property
    def tasks(self):
        """
        :return: dictionary of tasks and task_info by name.
        :rtype: dict
        """

        return self._tasks.copy()

    @tasks.setter
    def tasks(self, value):
        """
        Change of local dictionary of (task, task_info).

        :param dict value: new dictionary of tasks to use.
        """

        self._tasks = value

    def get_task(self, _id):
        """
        Get task registered at input _id.

        :param str _id: task identifier.
        :return: (task object, task info) related to input _id.
        :rtype: tuple
        :raises: KeyError if no task corresponds to _id. ImportError if task is
            impossible to load from runtime.
        """
        # if task is not already registered
        if _id not in self.tasks:
            # if file exists in DB
            _file = self[TaskManager.FILE_STORAGE].get(_id=_id)
            if _file is not None:
                # copy _file in self task directory
                copy(_file, self.task_directory)
                # import and reload the module
                module = import_module(_id)
                reload(module)

            # if task info exists in DB
            task_info = self[TaskManager.STORAGE].get_elements(ids=_id)

            # raises automatically an ImportError if task is not in runtime
            task = lookup(_id)

            # save task and task info in self tasks
            self.tasks[_id] = task, task_info

        # throws automatically a KeyError if _id is not in self.tasks
        result = self.tasks[_id]

        return result

    def set_task(self, _id, task=None, task_info=None, _file=None):
        """Change of task.

        :param str _id: task id.
        :param task: new task to use.
        :param dict task_info: task kwargs.
        :param _file: file where save the task
        :param callable task
        """

        # save task file if necessary
        if _file is not None:
            # in storage
            self[TaskManager.FILE_STORAGE].put(_id, _file)
            # and in directory
            copy(_file, self.task_directory)
            # and reload the module
            import_module(_file)
            reload(_file)

        # save task info
        self[TaskManager.FILE_STORAGE].put_element(_id=_id, element=task_info)

        # get task if None
        if task is None:
            task = lookup(_id)

        # save task in self tasks
        self.tasks[_id] = (task, task_info)
