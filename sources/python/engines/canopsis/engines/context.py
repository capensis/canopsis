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
        hostgroups = event.get('hostgroups', [])
        servicegroups = event.get('servicegroups', [])
        source_type = event['source_type']

        # get a copy of event
        _event = event.copy()

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

        # get related entity
        entity = self.context.get_entity(
            _event, from_db=True, create_if_not_exists=False
        )

        # get hostgroups
        entity['hostgroups'] = hostgroups
        # and service groups
        # create an entity status which is a component or a resource
        if source_type == 'resource':
            entity['servicegroups'] = servicegroups

        # create an entity status which is a component or a resource
        if source_type == 'resource':
            context['component'] = _event['component']
            entity['servicegroups'] = servicegroups

        if source_type not in ['resource', 'component']:
            self.logger.warning('source_type unknown %s' % source_type)

        # set mCrit and mWarn
        entity['mCrit'] = _event.get(mCrit, None)
        entity['mWarn'] = _event.get(mWarn, None)

        context, name = self.context.get_entity_context_and_name(entity)

        # put the status entity in the context
        self.put(_type=source_type, entity=entity, ctx=context)

        # udpdate context information with resource and component
        if source_type == 'resource':
            context['resource'] = name
        else:
            context['component'] = name

        # remove type from context because type will be metric
        del context['type']

        # add perf data (may be done in the engine perfdata)
        for perfdata in event.get('perf_data_array', []):
            perfdata_entity = entity.copy()
            name = perfdata['metric']
            perfdata_entity[Context.NAME] = name
            perfdata_entity['internal'] = perfdata['metric'].startswith('cps')
            self.put(_type='metric', entity=perfdata_entity, ctx=context)

        return event
