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
from canopsis.sla.core import Sla

import json

STATE_CRITICAL = 3
STATE_MAJOR = 2
STATE_MINOR = 1




class Watcher(MiddlewareRegistry):
    """Watcher"""

    OBJECT_STORAGE = ''
    ALERTS_STORAGE = ''

    def __init__(self, *args, **kwargs):
        """__init__

        :param *args:
        :param **kwargs:
        """
        super(Watcher, self).__init__(*args, **kwargs)
        self.object_storage = Middleware.get_middleware_by_uri(
            'storage-default://',
            table='object'
        )
        self[Watcher.OBJECT_STORAGE] = self.object_storage
        alerts_storage = Middleware.get_middleware_by_uri(
            'mongodb-periodical-alarm://'
        )
        self[Watcher.ALERTS_STORAGE] = alerts_storage

        self.sla_storage = Middleware.get_middleware_by_uri(
            'storage-default-sla://'
        )

        self.context_graph = ContextGraph()

    def create_watcher(self, body):
        """
            create watcher entity in context and link it to entities

            :param dict body: watcher conf
        """
        watcher_id = 'watcher-{0}'.format(
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
            id=watcher_id,
            name=body['display_name'],
            etype='watcher',
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
                '_id': watcher_id,
                'states': [0, 0, 0, 0, 0]
            }
        )

        self.calcul_state(watcher_id)

    def delete_watcher(self, watcher_id):
        """
            Delete watcher and disable entities linked to the watcher in context

            :param string watcher_id: watcher_id
        """
        object_watcher = list(self.object_storage._backend.find(
            {'_id': watcher_id}
        ))[0]
        watcher_entity = self.context_graph.get_entities_by_id(
            'watcher-{0}'.format(object_watcher['display_name'])
        )[0]
        watcher_entity['infos']['enabled'] = False
        self.context_graph.update_entity(watcher_entity)
        self.sla_storage.remove_elements(ids=[watcher_id])


    def calcul_state(self, watcher_id):
        """
            Compute state

            send an event watcher with the new state of the watcher

            :param watcher_id: watcher id
        """
        self.logger.debug('calcul')
        watcher_entity = self.context_graph.get_entities(
            query={'_id': watcher_id}
        )[0]

        entities = watcher_entity['depends']
        display_name = watcher_entity['name']

        alarm_list = list(self[Watcher.ALERTS_STORAGE]._backend.find(
            {'d': {'$in': entities}}
        ))
        states = []
        for alarm in alarm_list:
            if alarm['v']['resolved'] == None and alarm['d'] in entities:
                # add here a check of pebehavior to take into account or not
                # the alarm's state
                states.append(alarm['v']['state']['val'])

        nb_entities = len(entities)
        nb_crit = states.count(STATE_CRITICAL)
        nb_major = states.count(STATE_MAJOR)
        nb_minor = states.count(STATE_MINOR)
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

        self.debug(output)

        if computed_state != watcher_entity['infos']['state']:
            self.logger.critical('update entity')
            watcher_entity['infos']['state'] = computed_state
            self.context_graph.update_entity(watcher_entity)

        self.publish_event(display_name, computed_state, output)

    def worst_state(self, nb_crit, nb_major, nb_minor):
        """Worst state

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
        """publish_event

        Publish an event watcher on amqp

        :param display_name: watcher display_name
        :param computed_state: watcher state
        :param output: watcher output
        """
        event = forger(
            connector="canopsis",
            connector_name="engine",
            event_type="watcher",
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
        """alarm_changed

        Launch a computation of a watcher state

        :param alarm_id: alarm id
        """
        watchers = self.context_graph.get_entities(query={'type':'watcher'})
        for i in watchers:
            if alarm_id in i['depends']:
                self.calcul_state(i['_id'])

    def sla_compute(self, watcher_id, state):
        """sla_calcul
            launch the sla calcul

        :param watcher_id: watcher id
        :param state: watcher state
        """
        sla_tab = list(self.sla_storage.get_elements(query={'_id': watcher_id}))[0]
        sla_tab['states'][state] = sla_tab['states'][state] + 1

        self.sla_storage.put_element(sla_tab)

        watcher_conf = list(self.object_storage.get_elements(
            query={'_id':watcher_id}
        ))[0]


        sla = Sla(
            self.object_storage,
            'test/de/rk/on/verra/plus/tard',
            watcher_conf['sla_output_tpl'],
            watcher_conf['sla_timewindow'],
            watcher_conf['sla_warning'],
            watcher_conf['alert_level'],
            watcher_conf['display_name']
        )

        # self.logger.critical('{0}'.format((
        #     sla_tab['states']/
        #     (sla_tab['states'][1] +
        #      sla_tab['states'][2] +
        #      sla_tab['states'][3]))))

    def compute_slas(self):
        """
            launch the sla calcul for each watchers
        """
        watcher_list = self.context_graph.get_entities(
            query={'type': 'watcher', 'infos.enabled': True}
        )
        for watcher in watcher_list:
            self.sla_compute(
                watcher['_id'],
                watcher['infos']['state']
            )
