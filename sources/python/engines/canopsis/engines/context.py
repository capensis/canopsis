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
from __future__ import unicode_literals

from canopsis.engines.core import Engine
from canopsis.context.manager import Context
try:
    from threading import Lock
except ImportError:
    from dummy_threading import Lock

"""
TODO: sla
from canopsis.middleware.core import Middleware
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

        .. code-block:: python

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

    def work(self, event, *args, **kwargs):
        mCrit = 'PROC_CRITICAL'
        mWarn = 'PROC_WARNING'

        """
        TODO: sla

        .. code-block::

            if self.sla:
                mCrit = self.sla.data['mCrit']
                mWarn = self.sla.data['mWarn']
        """

        context = {}

        # Get event informations
        hostgroups = event.get('hostgroups', [])
        servicegroups = event.get('servicegroups', [])
        component = event.get('component')
        resource = event.get('resource')
        # quick fix when an event has an empty resource
        if 'resource' in event and not resource:
            del event['resource']

        # get a copy of event
        _event = event.copy()

        # add hostgroups
        for hostgroup in hostgroups:
            hostgroup_data = {
                Context.NAME: hostgroup
            }
            self.context.put(
                _type='hostgroup', entity=hostgroup_data, cache=True
            )
        # add servicegroups
        for servicegroup in servicegroups:
            servgroup_data = {
                Context.NAME: servicegroup
            }
            self.context.put(
                _type='servicegroup', entity=servgroup_data, cache=True
            )

        # get related entity
        encoded_event = {}
        for k, v in _event.items():
            try:
                k = k.encode('utf-8')
            except:
                pass
            try:
                v = v.encode('utf-8')
            except:
                pass
            encoded_event[k] = v

        entity = self.context.get_entity(
            encoded_event, from_db=True, create_if_not_exists=True, cache=True
        )

        # set service groups and hostgroups
        if resource:
            context['component'] = component
            entity['servicegroups'] = servicegroups
        entity['hostgroups'] = hostgroups

        # set mCrit and mWarn
        entity['mCrit'] = _event.get(mCrit, None)
        entity['mWarn'] = _event.get(mWarn, None)

        context, name = self.context.get_entity_context_and_name(entity)

        if 'resource' in context and not context['resource']:
            del context['resource']
        if 'resource' in entity and not entity['resource']:
            del entity['resource']

        # put the status entity in the context
        self.context.put(
            _type=entity[Context.TYPE], entity=entity, context=context,
            cache=True
        )

        # udpdate context information with resource and component
        if resource:
            context['resource'] = name
        else:
            context['component'] = name

        # remove type from context because type will be metric
        del context[Context.TYPE]

        # add perf data (may be done in the engine perfdata)
        for perfdata in event.get('perf_data_array', []):
            perfdata_entity = {}
            name = perfdata['metric']
            perfdata_entity[Context.NAME] = name
            perfdata_entity['internal'] = perfdata['metric'].startswith('cps')
            self.context.put(
                _type='metric', entity=perfdata_entity, context=context,
                cache=True
            )

        return event
