# -*- coding: UTF-8 -*-
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

from canopsis.context.manager import Context


class engine(Engine):
    etype = 'entities2'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.context = Context()

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
        connector = event['connector']
        connector_name = event['connector_name']
        component = event['component']
        resource = event.get('resource', None)
        hostgroups = event.get('hostgroups', [])
        servicegroups = event.get('servicegroups', [])
        source_type = event['source_type']
        event_type = event['event_type']

        # add connector
        entity = {
            'name': connector
        }
        self.context.put(
            _type='connector', entity=entity)

        # add connector_name
        entity['name'] = connector
        entity['connector'] = connector
        self.context.put(
            _type='connector_name', entity=entity)

        # add status entity which is a component or a resource
        entity['name'] = 'component'
        entity['connector_name'] = connector_name
        status_entity = entity.copy()
        status_entity['hostgroups'] = hostgroups

        is_status_entity = False

        # create an entity status which is a component or a resource
        if source_type == 'component':
            is_status_entity = True

        elif source_type == 'resource':
            # add component
            self.context.put(_type='component', entity=status_entity)
            is_status_entity = True
            status_entity['component'] = component
            status_entity['name'] = resource
            status_entity['servicegroups'] = servicegroups

        else:
            self.logger.warning('source_type unknown %s' % source_type)

        if is_status_entity:
            # add status entity
            status_entity['mCrit'] = event.get(mCrit, None)
            status_entity['mWarn'] = event.get(mWarn, None)
            status_entity['state'] = event['state']
            status_entity['state_type'] = event['state_type']
            self.context.put(_type=source_type, entity=status_entity)

        # add hostgroups
        for hostgroup in hostgroups:
            hostgroup_data = {
                'name': hostgroup
            }
            self.put(_type='hostgroup', entity=hostgroup_data)

        # add servicegroups
        for servicegroup in servicegroups:
            servicegroup_data = {
                'name': servicegroup
            }
            self.put(_type='servicegroup', entity=servicegroup_data)

        # add authored entity data (downtime, ack, metric, etc.)
        authored_data = entity.copy()
        if resource is not None:
            authored_data['resource'] = resource
            authored_data['author'] = event['author']
            authored_data['comment'] = event['comment']

        if event_type == 'ack':
            authored_data['timestamp'] = event['timestamp']

        elif event_type == 'downtime':
            authored_data['id'] = event['downtime_id']
            authored_data['start'] = event['start']
            authored_data['end'] = event['end']
            authored_data['duration'] = event['duration']
            authored_data['fixed'] = event['fixed']
            authored_data['entry'] = event['entry']

        authored_data['name'] = event_type

        self.context.put(_type=event_type, entity=authored_data)

        # add perf data
        for perfdata in event['perf_data_array']:
            perfdata_entity = entity.copy()
            perfdata_entity['name'] = perfdata['metric']
            perfdata_entity['internal'] = perfdata['metric'].startswith('cps')
            self.context.put(_type='metric', entity=perfdata_entity)

        return event
