# -*- coding: UTF-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.engines import Engine
from canopsis.old.account import Account
from canopsis.old.storage import get_storage

from md5 import new


class engine(Engine):
    etype = 'entities'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.account = Account(user='root', group='root')
        self.storage = get_storage(
            namespace='entities', logging_level=self.logging_level,
            account=self.account)
        self.backend = self.storage.get_backend()

        self.sla = None
        self.beat()

    def update(self, doc, hint):
        if not self.backend.find(doc).hint(hint).limit(-1).count():
            self.backend.save(doc)

    def beat(self):
        cursor = self.storage.get_backend('object').find({
            'crecord_type': 'sla',
            'objclass': 'macro'
        }).hint([('crecord_type', 1)]).limit(-1)

        if cursor.count():
            self.sla = cursor[0]

    def work(self, event, *args, **kwargs):
        mCrit = 'PROC_CRITICAL'
        mWarn = 'PROC_WARNING'

        if self.sla:
            mCrit = self.sla.data['mCrit']
            mWarn = self.sla.data['mWarn']

        # Get event informations
        component = event['component']
        resource = event.get('resource', None)
        hostgroups = event.get('hostgroups', [])
        servicegroups = event.get('servicegroups', [])

        # Create Component entity
        doc_id = 'component.{0}'.format(component)
        doc = self.backend.find_one(doc_id)

        if not doc:
            doc = {
                '_id': doc_id,
                'type': 'component',
                'name': component
            }

        doc['hostgroups'] = hostgroups

        if event['source_type'] == 'component':
            doc['mCrit'] = event.get(mCrit, None)
            doc['mWarn'] = event.get(mWarn, None)

            doc['state'] = event['state']
            doc['state_type'] = event['state_type']

        self.update(doc, [('type', 1), ('name', 1)])

        # Create Resource entity
        if resource:
            doc = {
                '_id': 'resource.{0}.{1}'.format(component, resource),
                'type': 'resource',
                'name': resource,
                'component': component,
                'hostgroups': hostgroups,
                'servicegroups': servicegroups,
                'mCrit': event.get(mCrit, None),
                'mWarn': event.get(mWarn, None),

                'state': event['state'],
                'state_type': event['state_type']
            }

            self.update(doc, [('type', 1), ('component', 1), ('name', 1)])

        # Create Hostgroups entities
        for hostgroup in hostgroups:
            doc = {
                '_id': 'hostgroup.{0}'.format(hostgroup),
                'type': 'hostgroup',
                'name': hostgroup
            }

            self.update(doc, [('type', 1), ('name', 1)])

        # Create Servicegroups entities
        for servicegroup in servicegroups:
            doc = {
                '_id': 'servicegroup.{0}'.format(servicegroup),
                'type': 'servicegroup',
                'name': servicegroup
            }

            self.update(doc, [('type', 1), ('name', 1)])

        # Create Downtime entity
        if event['event_type'] == 'downtime':
            doc = {
                '_id': 'downtime.{0}.{1}.{2}'.format(
                    component, resource, event['downtime_id']),
                'type': 'downtime',
                'component': component,
                'resource': resource,
                'id': event['downtime_id'],

                'author': event['author'],
                'comment': event['output'],

                'start': event['start'],
                'end': event['end'],
                'duration': event['duration'],

                'fixed': event['fixed'],
                'entry': event['entry']
            }

            self.update(
                doc,
                [('type', 1), ('component', 1), ('resource', 1), ('id', 1)])

        # Create acknowledgement entity
        elif event['event_type'] == 'ack':
            doc = {
                '_id': 'ack.{0}.{1}.{2}'.format(
                    component, resource, event['timestamp']),
                'type': 'ack',
                'timestamp': event['timestamp'],
                'component': component,
                'resource': resource,

                'author': event['author'],
                'comment': event['output'],
            }

            self.update(doc, [('type', 1), ('component', 1), ('resource', 1)])

        # Create metrics entities
        for perfdata in event.get('perf_data_array', []):
            nodeid = new()

            nodeid.update(component.encode('ascii', 'ignore'))

            if resource:
                nodeid.update(resource.encode('ascii', 'ignore'))

            nodeid.update(perfdata['metric'])
            nodeid = nodeid.hexdigest()

            doc = {
                '_id': 'metric.{0}'.format(nodeid),
                'type': 'metric',
                'component': component,
                'resource': resource,
                'name': perfdata['metric'],
                'nodeid': nodeid,
                'internal': perfdata['metric'].startswith('cps_')
            }

            self.update(doc, [('type', 1), ('nodeid', 1)])

        return event
