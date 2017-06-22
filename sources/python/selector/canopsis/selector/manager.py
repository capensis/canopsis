# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.middleware.core import Middleware

from canopsis.context_graph.manager import ContextGraph
from canopsis.context_graph.process import create_entity

from canopsis.engines.core import publish
from canopsis.event import forger, get_routingkey
from canopsis.old.rabbitmq import Amqp
from canopsis.sla.core import Sla

import json


class Selector(MiddlewareRegistry):
    """Selector"""

    OBJECT_STORAGE = ''
    ALERTS_STORAGE = ''

    def __init__(self, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
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

        self.sla_storage = Middleware.get_middleware_by_uri(
            'storage-default-sla://'
        )

        self.context_graph = ContextGraph()

    def create_selector(self, body):
        """
        Create selector entity in context and link to entities.

        :param dict body: selector conf
        """
        selector_id = 'selector-{0}'.format(
            body['display_name']

        )
        depends_list = self.context_graph.get_entities(
            query=json.loads(body['mfilter']),
            projection={'_id': 1}
        )
        depend_list = []
        for entity_id in depends_list:
            depend_list.append(entity_id['_id'])

        entity = create_entity(
            id=selector_id,
            name=body['display_name'],
            etype='selector',
            impact=[],
            depends=depend_list,
            infos={
                'mfilter': body['mfilter'],
                'enabled': True,
                'state': 0
            }
        )
        self.context_graph.create_entity(entity)

        self.sla_storage.put_element(
            element={
                '_id': selector_id,
                'states': [0, 0, 0, 0, 0]
            }
        )

        self.calcul_state(selector_id)

    def delete_selector(self, selector_id):
        """
        Delete_selector & disable selector entity in context.

        :param string selector_id: selector_id
        """
        object_selector = list(self.object_storage._backend.find(
            {'_id': selector_id}
        ))[0]
        selector_entity = self.context_graph.get_entities_by_id(
            'selector-{0}'.format(object_selector['display_name'])
        )[0]
        selector_entity['infos']['enabled'] = False
        self.context_graph.update_entity(selector_entity)
        self.sla_storage.remove_elements(ids=[selector_id])

    def calcul_state(self, selector_id):
        """
        Send an event selector with the new state of the selector.

        :param selector_id: selector id
        """
        self.logger.critical('calcul')
        selector_entity = self.context_graph.get_entities(
            query={'_id': selector_id}
        )[0]

        entities = selector_entity['depends']
        display_name = selector_entity['name']

        alarm_list = list(self[Selector.ALERTS_STORAGE]._backend.find(
            {'d': {'$in': entities}}
        ))
        states = []
        for alarm in alarm_list:
            if alarm['v']['resolved'] is None and alarm['d'] in entities:
                # add here a check of pebehavior to take into account or not
                # the alarm's state
                states.append(alarm['v']['state']['val'])

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

        if computed_state != selector_entity['infos']['state']:
            self.logger.critical('update entity')
            selector_entity['infos']['state'] = computed_state
            self.context_graph.update_entity(selector_entity)

        self.publish_event(display_name, computed_state, output)

    def worst_state(self, nb_crit, nb_major, nb_minor):
        """Calculate the worst state.

        :param int nb_crit: critical number
        :param int nb_major: major number
        :param int nb_minor: minor number
        :return int state: return the worst state
        """

        if nb_crit > 0:
            return 3
        elif nb_major > 0:
            return 2
        elif nb_minor > 0:
            return 1
        else:
            return 0

    def publish_event(self, display_name, computed_state, output):
        """
        Publish an event selector on amqp.

        :param display_name: selector display_name
        :param computed_state: selector state
        :param output: selector output
        """
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
        #self.logger.critical('published {0}'.format(event))

    def alarm_changed(self, alarm_id):
        """
        Launch a caculation of a selector state.

        :param alarm_id: alarm id
        """
        selectors = self.context_graph.get_entities(query={'type':'selector'})
        for i in selectors:
            if alarm_id in i['depends']:
                self.calcul_state(i['_id'])

    def sla_compute(self, selector_id, state):
        """
        Launch the sla calcul.

        :param selector_id: selector id
        :param state: selector state
        """
        sla_tab = list(self.sla_storage.get_elements(query={'_id': selector_id}))[0]
        sla_tab['states'][state] = sla_tab['states'][state] + 1

        self.sla_storage.put_element(sla_tab)

        selector_conf = list(self.object_storage.get_elements(
            query={'_id': selector_id}
        ))[0]

        sla = Sla(
            self.object_storage,
            'test/de/rk/on/verra/plus/tard',
            selector_conf['sla_output_tpl'],
            selector_conf['sla_timewindow'],
            selector_conf['sla_warning'],
            selector_conf['alert_level'],
            selector_conf['display_name']
        )
        """
        self.logger.critical('{0}'.format((
            sla_tab['states']/
            (sla_tab['states'][1] +
             sla_tab['states'][2] +
             sla_tab['states'][3]))))
        """

    def compute_slas(self):
        """
        Launch the sla calcul for each selectors.
        """
        selector_list = self.context_graph.get_entities(
            query={'type': 'selector', 'infos.enabled': True}
        )
        for selector in selector_list:
            self.sla_compute(
                selector['_id'],
                selector['infos']['state']
            )
