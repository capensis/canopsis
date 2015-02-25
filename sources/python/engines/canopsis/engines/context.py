# -*- coding: UTF-8 -*-
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

from canopsis.engines import Engine
from canopsis.context.manager import Context
try:
    from threading import Lock
except ImportError:
    from dummy_threading import Lock

"""
TODO: sla
from canopsis.middleware import Middleware
"""


class engine(Engine):
    etype = 'context'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        # get a context
        self.context = Context()
        """
        TODO: sla
        # get a storage for sla macro
        #self.storage = Middleware.get_middleware(
        #    protocol='storage', data_scope='global')

        #self.sla = None
        """

        self.entities_by_entity_ids = {}
        self.lock = Lock()
        self.beat()

    def beat(self):
        """
        TODO: sla
        sla = self.storage.find_elements(request={
            'crecord_type': 'sla',
            'objclass': 'macro'
        })

        if sla:
            self.sla = sla[0]
        """

        self.lock.acquire()
        entities_by_entity_ids = self.entities_by_entity_ids.copy()
        self.entities_by_entity_ids = {}
        self.lock.release()

        entities = self.context[Context.CTX_STORAGE].get_elements(
            ids=entities_by_entity_ids.keys()
        )

        for entity in entities:
            del entities_by_entity_ids[entity['_id']]

        if entities_by_entity_ids:
            context = self.context
            for entity_id in entities_by_entity_ids:
                _type, entity, ctx = entities_by_entity_ids[entity_id]
                context.put(
                    _type=_type, entity=entity, context=ctx
                )

    def put(self, _type, entity, ctx=None):

        context = self.context
        full_entity = entity.copy()
        full_entity[Context.TYPE] = _type
        if ctx is not None:
            full_entity.update(ctx)
        eid = context.get_entity_id(full_entity)
        self.lock.acquire()
        self.entities_by_entity_ids[eid] = _type, entity, ctx
        self.lock.release()

    def work(self, event, *args, **kwargs):
        mCrit = 'PROC_CRITICAL'
        mWarn = 'PROC_WARNING'

        """
        TODO: sla
        if self.sla:
            mCrit = self.sla.data['mCrit']
            mWarn = self.sla.data['mWarn']
        """

        context = {}

        # Get event informations
        connector = event['connector']
        connector_name = event['connector_name']
        component = event['component']
        resource = event.get('resource', None)
        hostgroups = event.get('hostgroups', [])
        servicegroups = event.get('servicegroups', [])

        # add connector
        entity = {}
        # add connector_name
        context['connector'] = connector
        context['connector_name'] = connector_name

        # add status entity which is a component or a resource
        entity[Context.NAME] = component if resource is None else resource

        status_entity = entity.copy()
        status_entity['hostgroups'] = hostgroups

        is_status_entity = False
        source_type = event['source_type']

        # create an entity status which is a component or a resource
        if source_type == 'resource':
            status_entity['servicegroups'] = servicegroups

        if source_type not in ['resource', 'component']:
            self.logger.warning('source_type unknown %s' % source_type)

        is_status_entity = True

        if is_status_entity:
            # add status entity
            status_entity['mCrit'] = event.get(mCrit, None)
            status_entity['mWarn'] = event.get(mWarn, None)
            status_entity['state'] = event['state']
            status_entity['state_type'] = event['state_type']

        # add hostgroups
        for hostgroup in hostgroups:
            hostgroup_data = {
                Context.NAME: hostgroup
            }
            self.put(_type='hostgroup', entity=hostgroup_data)

        # add servicegroups
        for servicegroup in servicegroups:
            servgroup_data = {
                Context.NAME: servicegroup
            }
            self.put(_type='servicegroup', entity=servgroup_data)

        # put the status entity in the context
        self.put(_type=source_type, entity=status_entity, ctx=context)

        # udpdate context information with resource and component
        context['component'] = component
        if resource is not None:
            context['resource'] = resource

        authored_data = entity.copy()
        event_type = event['event_type']

        if 'author' in event:
            # add authored entity data (downtime, ack, etc.)
            authored_data['author'] = event['author']
            authored_data['comment'] = event.get('output', None)

            if authored_data['comment'] is None:
                del authored_data['comment']

            if event_type == 'ack':
                authored_data['timestamp'] = event['timestamp']
                authored_data[Context.NAME] = str(event['timestamp'])

            elif event_type == 'downtime':
                authored_data['downtime_id'] = event['downtime_id']
                authored_data['start'] = event['start']
                authored_data['end'] = event['end']
                authored_data['duration'] = event['duration']
                authored_data['fixed'] = event['fixed']
                authored_data['entry'] = event['entry']
                authored_data[Context.NAME] = event['rk']

        self.put(_type=event_type, entity=authored_data, ctx=context)

        # add perf data
        for perfdata in event.get('perf_data_array', []):
            perfdata_entity = entity.copy()
            name = perfdata['metric']
            perfdata_entity[Context.NAME] = name
            perfdata_entity['internal'] = perfdata['metric'].startswith('cps')
            self.put(_type='metric', entity=perfdata_entity, ctx=context)

        return event
