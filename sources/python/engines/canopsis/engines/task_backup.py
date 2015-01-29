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

from canopsis.engines import TaskHandler

import shutil

import os

import ubik.core as ubik_api

from subprocess import Popen


class engine(TaskHandler):
    etype = 'task_backup'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.backup_path = os.path.expanduser('~/var/backups')
        self.mongo_archive = 'backup_mongo'
        self.config_archive = 'backup_config'

        self.mongo_dir = os.path.join(self.backup_path, self.mongo_archive)
        self.mongo_path = '{0}.tar.gz'.format(self.mongo_dir)

        self.config_dir = os.path.join(self.backup_path, self.config_archive)
        self.config_path = '{0}.tar.gz'.format(self.config_dir)

        self.etc_dir = os.path.expanduser('~/etc')
        self.pkg_file = os.path.expanduser('~/etc/.packages')

    def handle_task(self, job):
        action = job.get('action', 'all')
        hostname = job.get('hostname', 'localhost')

        try:
            if action == 'mongo':
                self.backup_mongo(hostname)

            elif action == 'config':
                self.backup_config()

            else:  # action == 'all'
                self.backup_mongo()
                self.backup_config()

        except IOError as err:
            return (
                2,
                'Backup failed: {0}'.format(err)
            )

        return (0, 'Backup done')

    def backup_mongo(self, hostname):
        if os.path.exists(self.mongo_dir):
            self.logger.debug('RMDIR: {0}'.format(self.mongo_dir))
            shutil.rmtree(self.mongo_dir)

        if os.path.exists(self.mongo_path):
            self.logger.debug('RMBAK: {0}'.format(self.mongo_path))
            os.remove(self.mongo_path)

        self.logger.debug('MKDIR: {0}'.format(self.mongo_dir))
        os.makedirs(self.mongo_dir)

        cmd = 'mongodump --host {0} --out {1}/'.format(
            hostname, self.mongo_dir)
        self.logger.debug('EXEC: {0}'.format(cmd))
        dumpout = Popen(cmd, shell=True)
        dumpout.wait()

        self.logger.debug('MKTAR: {0}'.format(self.mongo_path))
        shutil.make_archive(
            base_name=os.path.join(self.backup_path, self.mongo_archive),
            root_dir=self.backup_path,
            format='gztar',
            base_dir=self.mongo_archive,
            logger=self.logger
        )

        self.logger.debug('RMDIR: {0}'.format(self.mongo_dir))
        shutil.rmtree(self.mongo_dir)

        if not os.path.exists(self.mongo_path):
            raise IOError(
                'Archive file "{0}" not found'.format(self.mongo_path))

    def backup_config(self):
        if os.path.exists(self.config_dir):
            self.logger.debug('RMDIR: {0}'.format(self.config_dir))
            shutil.rmtree(self.config_dir)

        if os.path.exists(self.config_path):
            self.logger.debug('RMBAK: {0}'.format(self.config_path))
            os.remove(self.config_path)

        self.logger.debug('LSPKG: {0}'.format(self.pkg_file))
        lines = []
        for package in ubik_api.db.get_installed():
            lines.append(package.name)
            lines.append('\n')

        del lines[-1]

        with open(self.pkg_file, 'w') as f:
            f.writelines(lines)

        self.logger.debug('COPYCONF: {0}'.format(self.config_dir))
        shutil.copytree(self.etc_dir, self.config_dir)

        self.logger.debug('MKTAR: {0}'.format(self.config_path))
        shutil.make_archive(self.config_dir, 'gztar', self.config_dir)

        self.logger.debug('RMDIR: {0}'.format(self.config_dir))
        shutil.rmtree(self.config_dir)

        if not os.path.exists(self.config_path):
            raise IOError('Archive "{0}" not found'.format(self.config_path))
