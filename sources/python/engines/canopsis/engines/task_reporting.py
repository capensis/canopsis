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

from canopsis.old.storage import get_storage
from canopsis.old.record import Record
from canopsis.old.account import Account
from canopsis.old.file import File

from wkhtmltopdf.wrapper import Wrapper

import hashlib
import datetime
import time
import json
import os


class engine(TaskHandler):
    etype = 'task_reporting'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

    def handle_task(self, job):
        now = int(time.time())

        filename = job.get('filename', None)
        viewname = job.get('viewname', None)

        start = job.get('start', -1)
        stop = job.get('stop', now)
        subset_selection = job.get('subset_selection', None)
        interval = job.get('interval', None)

        account = job.get('account', None)
        owner = job.get('owner', None)
        mail = job.get('mail', None)
        orientation = job.get('orientation', 'Portrait')
        pagesize = job.get('pagesize', 'A4')

        # Verify interval

        if interval and start < 0:
            start = stop - interval

        raccount = Account(user='root', group='root', mail='root@localhost.local')
        rstorage = get_storage('object', account=raccount)

        taccount = rstorage.find_one(
            mfilter={
                '_id': 'account.{0}'.format(account)
            })

        if isinstance(taccount, Record):
            account = Account(taccount)

            self.logger.info('User {0} retrieved from database'.format(account.user))

        else:
            account = Account(mail='anonymous@localhost.local')

            self.logger.info('No user found in database, identified as anonymous')

        if owner:
            towner = rstorage.find_one(
                mfilter={
                    '_id': 'account.{0}'.format(owner)
                })

            if isinstance(towner, Record):
                owner = Account(towner)

        # Fetch view
        storage = get_storage('object', account=account)

        try:
            view = storage.get(viewname)

        except Exception:
            return (
                2,
                "Can't find view {0} with account {1}".format(
                    viewname, account.user)
            )

        self.logger.info(
            'Account {0} requested rendering the view {1}'.format(account.user, view.name))

        # Generate filename
        if not filename:
            dtstart = datetime.date.fromtimestamp(stop)

            if start > 0:
                dtend = datetime.date.fromtimestamp(start)

                filename = '{0}_From_{1}_To_{2}.pdf'.format(
                    view.name,
                    dtstart,
                    dtend
                )

            else:
                filename = '{0}_{1}.pdf'.format(view.name, dtend)

        filename = hashlib.md5(filename.encode('ascii', 'ignore')).hexdigest()

        # View options
        viewopts = view.data.get('view_options', {})

        if isinstance(viewopts, dict):
            orientation = viewopts.get('orientation', orientation)
            pagesize = viewopts.get('pagesize', pagesize)

        # Run report
        doc_id = self.report(
            filename,
            viewname,
            start,
            stop,
            subset_selection,
            account,
            owner,
            orientation,
            pagesize
        )

        if not doc_id:
            return (
                2,
                'Impossible to save report in GridFS'
            )

        # Send mail

        self.mail(
            job,
            mail,
            account,
            doc_id
        )

        return (0, 'View {0} rendered by {1}'.format(view.name, account.user))

    def report(
        self,
        filename, viewname, start, stop, subset_selection, account, owner,
        orientation, pagesize
    ):
        wrapper_cfgpath = os.path.expanduser('~/etc/wkhtmltopdf_wrapper.json')
        fpath = None

        with open(wrapper_cfgpath, 'r') as f:
            conf = json.loads(f.read())

            fpath = os.path.join(
                conf['report_dir'],
                filename
            )

        wrapper = Wrapper(
            filename,
            viewname,
            start,
            stop,
            subset_selection,
            account,
            wrapper_cfgpath,
            orientation=orientation,
            pagesize=pagesize
        )

        self.logger.debug('Rendering PDF: {0} -> {1}'.format(viewname, fpath))
        wrapper.run_report()

        doc_id = self.put_in_GridFS(fpath, filename, account, owner)
        os.remove(fpath)

        return doc_id

    def mail(self, job, mail, account, doc_id):
        if isinstance(mail, dict):
            mail['id'] = str(job['id'])
            mail['user'] = account.user
            mail['attachments'] = [str(doc_id)]

            self.run_task('mail', mail)

    def put_in_GridFS(self, filepath, filename, account, owner):
        storage = get_storage('files', account)
        report = File(storage=storage)
        report.put_file(filepath, filename, content_type='application/pdf')

        if owner:
            report.chown(owner)

        doc_id = storage.put(report)

        if not report.check(storage):
            return False

        return doc_id
