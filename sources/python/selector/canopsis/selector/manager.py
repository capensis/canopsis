# -*- coding: utf-8 -*- 
from __future__ import unicode_literals

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.middleware.core import Middleware

from canopsis.context_graph.manager import ContextGraph
from canopsis.context_graph.process import create_entity

from canopsis.engines.core import publish
from canopsis.event import forger, get_routingkey
from canopsis.old.rabbitmq import Amqp

import json


class Selector(MiddlewareRegistry):

    OBJECT_STORAGE = ''
    ALERTS_STORAGE = ''

    def __init__(self, *args, **kwargs):
        super(Selector, self).__init__(*args, **kwargs)
        self.object_storage = Middleware.get_middleware_by_uri(
            'storage-default://',
            table='object'
        )
        self[Selector.OBJECT_STORAGE] = self.object_storage
        alerts_storage = Middleware.get_middleware_by_uri(
            'mongodb-periodical-alarm://'
        )
        self[Selector.ALERTS_STORAGE] = alerts_storage
        
        self.context_graph = ContextGraph()

    def create_selector(self, body):
        """
            create selector entity in context and link to entities
        """
        selector_id = '/selector/{0}'.format(
            body['display_name']

        )
        depends_list = self.context_graph.get_entities(
            query=json.loads(body['mfilter']),
            projection={'_id': 1}
        )
        d = []
        for i in depends_list:
            d.append(i['_id'])

        entity = create_entity(
            id=selector_id, 
            name=body['display_name'],
            etype='selector',
            impact=[],
            depends=d,
            infos={
                'mfilter': body['mfilter'],
                'enabled': True,
            }
        )
        self.context_graph.create_entity(entity)

        self.calcul_state(selector_id)

    def delete_selector(self, selector_id):
        """
            disable selector entity in context
        """
        object_selector = list(self.object_storage._backend.find(
            {'_id': selector_id}
        ))[0]
        selector_entity = self.context_graph.get_entities_by_id(
            '/selector/{0}'.format(object_selector['display_name'])
        )[0]
        selector_entity['infos']['enabled'] = False
        self.logger.critical(selector_entity)
        self.context_graph.update_entity(selector_entity)
        

    def calcul_state(self, selector_id):
        """
            send an event selector with the new state of the selector
        """
        s = self.context_graph.get_entities(
            query={'_id': selector_id},
            projection={'depends':1, 'name': 1}
        )[0]

        entities = s['depends']
        display_name = s['name']

        r = list(self[Selector.ALERTS_STORAGE]._backend.find(
            {'d': {'$in': entities}}
        ))
        states = []
        for i in r:
            if i['v']['resolved'] == None and i['d'] in entities:
                # add here a check of pebehavior to take into account or not
                # the alarm's state
                states.append(i['v']['state']['val']) 

        nb_entities = len(entities)
        nb_crit = states.count(3)
        nb_major = states.count(2)
        nb_minor = states.count(1)
        nb_ok = nb_entities - (nb_crit + nb_major + nb_minor)
        
        # here add selection for calculation method actually it's worst state 
        # by default and think to add pbehavior in tab
        computed_state = self.worst_state(nb_crit, nb_major, nb_minor)
        output = '{0} ok, {1} minor, {2} major, {3} critical'.format(
            nb_ok,
            nb_minor,
            nb_major,
            nb_crit
        )

        self.publish_event(display_name, computed_state, output)

    def worst_state(self, nb_crit, nb_major, nb_minor):
        if nb_crit > 0:
            return 3
        elif nb_major > 0:
            return 2
        elif nb_minor > 0:
            return 1
        else:
            return 0

    def publish_event(self, display_name, computed_state, output):
        event = forger(
            connector="canopsis",
            connector_name="engine",
            event_type="selector",
            source_type="component",
            component=display_name,
            state=computed_state,
            output=output,
            perf_data_array=[],
            display_name=display_name
        )

        rk = get_routingkey(event)
        amqp = Amqp()
        publish(event=event, publisher=amqp, rk=rk, logger=self.logger)
        self.logger.critical('published {0}'.format(event))

    def alarm_changed(self, alarm_id):
        selectors = self.context_graph.get_entities(query={'type':'selector'})
        for i in selectors:
            if alarm_id in i['depends']:
                self.calcul_state(i['_id'])

    def sla_calcul(self, selector_id):
        """
            launch the sla calcul
        """

