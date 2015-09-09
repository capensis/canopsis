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

from canopsis.engines.core import Engine, publish
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis import schema

from dateutil.rrule import rrulestr

from datetime import datetime

from time import time


class engine(Engine):
    etype = 'scheduler'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        account = Account(user='root', group='root')
        self.storage = get_storage('object', account=account)

    def work(self, job, *args, **kwargs):
        if schema.validate(job, 'crecord.job'):
            self.do_job(job)

        else:
            self.logger.error('Invalid job: {0}'.format(job))

    def pre_run(self):
        self.beat()

    def beat(self):
        self.logger.info('Reload jobs')

        now = int(time())
        prev = now - self.beat_interval

        jobs = self.storage.find({'$and': [
            {'crecord_type': 'job'},
            {'$or': [
                {'last_execution': {'$lte': prev}},
                {'last_execution': None},
            ]}
        ]})

        for job in jobs:
            job = job.dump()

            self.logger.debug('Job: {0}'.format(job))

            if job['last_execution'] <= 0:
                self.do_job(job)

            else:
                jobStart = datetime.fromtimestamp(job['start'])

                dtstart = datetime.fromtimestamp(prev)
                dtend = datetime.fromtimestamp(now)

                occurences = list(
                    rrulestr(job['rrule'], dtstart=jobStart).between(
                        dtstart, dtend))

                if len(occurences) > 0:
                    self.do_job(job)

    def do_job(self, job):
        self.logger.info('Execute job: {0}'.format(job))

        job['params']['jobid'] = job['_id']
        job['params']['jobctx'] = job.get('context', {})

        publish(
            publisher=self.amqp, event=job['params'],
            rk='task_{0}'.format(job['task'][4:]), exchange='amq.direct'
        )

        now = int(time())

        self.storage.get_backend().update({'$and': [
            {'_id': job['_id']},
            {'$or': [
                {'last_execution': {'$lte': now}},
                {'last_execution': None},
            ]}
        ]}, {
            '$set': {
                'last_execution': now
            }
        })
