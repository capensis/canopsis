# -*- coding: utf-8 -*-

"""
Manager for watcher.
"""

from __future__ import unicode_literals
import time
import json

from canopsis.check import Check
from canopsis.context_graph.manager import ContextGraph
from canopsis.event import forger
from canopsis.logger import Logger
from canopsis.common.middleware import Middleware
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.common.amqp import AmqpPublisher
from canopsis.common.amqp import get_default_connection as \
    get_default_amqp_conn

LOG_PATH = 'var/log/watcher.log'


class Watcher:
    """Watcher class"""

    def __init__(self, amqp_pub=None):
        """
        :param amqp_pub canopsis.common.amqp.AmqpPublisher:
        """
        self.logger = Logger.get('watcher', LOG_PATH)
        self.watcher_storage = Middleware.get_middleware_by_uri(
            'mongodb-default-watcher://')
        self.alert_storage = Middleware.get_middleware_by_uri(
            'mongodb-periodical-alarm://')

        self.sla_storage = Middleware.get_middleware_by_uri(
            'storage-default-sla://')

        self.context_graph = ContextGraph(self.logger)
        self.pbehavior_manager = PBehaviorManager(
            *PBehaviorManager.provide_default_basics()
        )
        self.amqp_pub = amqp_pub
        if amqp_pub is None:
            self.amqp_pub = AmqpPublisher(get_default_amqp_conn())

    def get_watcher(self, watcher_id):
        """Retreive from database the watcher specified by is watcher id.

        :param str watcher_id: the watcher id
        :return dict: the wanted watcher. None, if no watcher match the
        watcher_id
        """
        watcher = self.context_graph.get_entities_by_id(watcher_id)

        try:
            return watcher[0]
        except IndexError:
            return None

    def create_watcher(self, body):
        """
        Create watcher entity in context and link to entities.

        :param dict body: watcher conf
        """
        watcher_id = body['_id']
        try:
            watcher_finder = json.loads(body['mfilter'])
        except ValueError:
            self.logger.error('can t decode mfilter')
            return None
        except KeyError:
            self.logger.error('no filter')
            return None

        depends_list = self.context_graph.get_entities(
            query=watcher_finder,
            projection={'_id': 1}
        )
        self.watcher_storage.put_element(body)

        depend_list = []
        for entity_id in depends_list:
            depend_list.append(entity_id['_id'])

        entity = ContextGraph.create_entity_dict(
            id=watcher_id,
            name=body['display_name'],
            etype='watcher',
            impact=[],
            depends=depend_list
        )

        # adding the fields specific to the Watcher entities
        entity['mfilter'] = body['mfilter']
        entity['state'] = 0

        try:
            self.context_graph.create_entity(entity)
        except ValueError:
            self.context_graph.update_entity(entity)

        self.compute_state(watcher_id)

        return True  # TODO: return really something

    def update_watcher(self, watcher_id, updated_field):
        """Update the watcher specified by is watcher id with updated_field.

        Raise a ValueError, if the watcher_id do not match any entity.

        :param str watcher_id: the watcher_id of the watcher to update
        :param dict updated_field: the fields to update
        :returns: the updated Watcher
        :rtype: <Watcher>
        """

        watcher = self.get_watcher(watcher_id)

        if watcher is None:
            raise ValueError("No watcher found for the following"
                             " id: {}".format(watcher_id))

        if "mfilter" in watcher.keys() and "mfilter" in updated_field.keys():
            if updated_field['mfilter'] != watcher['mfilter']:
                watcher['mfilter'] = updated_field['mfilter']

                query = json.loads(updated_field['mfilter'])
                entities = self.context_graph.get_entities(
                    query=query, projection={'_id': 1})

                watcher["depends"] = [entity["_id"] for entity in entities]

        for key in updated_field:

            if key == "infos":  # update fields inside infos
                for info_key in updated_field["infos"]:
                    watcher["infos"][info_key] = updated_field["infos"][
                        info_key]

            watcher[key] = updated_field[key]

        self.context_graph.update_entity(watcher)

    def delete_watcher(self, watcher_id):
        """
        Delete watcher & disable watcher entity in context.

        :param string watcher_id: watcher_id
        :returns: the mongodb dict response
        """
        self.context_graph.delete_entity(watcher_id)

        self.sla_storage.remove_elements(ids=[watcher_id])

        return self.watcher_storage.remove_elements(ids=[watcher_id])

    def alarm_changed(self, alarm_id):
        """
        Launch a computation of a watcher state.

        :param alarm_id: alarm id
        """
        watchers = self.context_graph.get_entities(query={'type': 'watcher'})
        for i in watchers:
            if alarm_id in i['depends']:
                self.compute_state(i['_id'])

    def compute_watchers(self):
        """
        Compute all watchers states.
        """
        watchers = list(self.watcher_storage.get_elements(query={}))
        for watcher in watchers:
            self.compute_state(watcher['_id'])

    def compute_state(self, watcher_id):
        """
        Send an event watcher with the new state of the watcher.

        :param watcher_id: watcher id
        """
        try:
            watcher_entity = self.context_graph.get_entities(
                query={'_id': watcher_id})[0]
        except IndexError:
            return None

        entities = watcher_entity['depends']

        query = {"_id": {"$in": entities},
                 "enabled": True}
        cursor = self.context_graph.get_entities(query=query,
                                                 projection={"_id": 1})

        entities = []
        for ent in cursor:
            entities.append(ent["_id"])

        display_name = watcher_entity['name']

        alarm_list = list(self.alert_storage._backend.find({
            '$and': [
                {'d': {'$in': entities}},
                {
                    '$or': [
                        {'v.resolved': None},
                        {'v.resolved': {'$exists': False}}
                    ]
                }
            ]
        }))
        states = []

        for alarm in alarm_list:
            pbh_alarm = self.pbehavior_manager.get_pbehaviors_by_eid(alarm['d'])

            active_pbh = []
            now = int(time.time())
            for pbh in pbh_alarm:
                if self.pbehavior_manager.check_active_pbehavior(now, pbh):
                    active_pbh.append(pbh)
            if len(active_pbh) == 0:
                states.append(alarm['v']['state']['val'])

        nb_entities = len(entities)
        nb_crit = states.count(Check.CRITICAL)
        nb_major = states.count(Check.MAJOR)
        nb_minor = states.count(Check.MINOR)
        nb_ok = nb_entities - (nb_crit + nb_major + nb_minor)

        # here add selection for calculation method actually it's worst state
        # by default and think to add pbehavior in tab
        computed_state = self.worst_state(nb_crit, nb_major, nb_minor)
        output = '{0} ok, {1} minor, {2} major, {3} critical'.format(
            nb_ok, nb_minor, nb_major, nb_crit)

        if computed_state != watcher_entity.get('state', None):
            watcher_entity['state'] = computed_state
            self.context_graph.update_entity(watcher_entity)

        self.publish_event(
            display_name,
            computed_state,
            output,
            watcher_entity['_id']
        )

    def compute_slas(self):
        """
        Launch the sla calcul for each watchers.
        """
        watcher_list = self.context_graph.get_entities(
            query={'type': 'watcher',
                   'infos.enabled': True})
        for watcher in watcher_list:
            self.sla_compute(watcher['_id'], watcher['infos']['state'])

    def publish_event(self, display_name, computed_state, output, _id):
        """
        Publish an event watcher on amqp.

        TODO: move that elsewhere (not specific to watchers)

        :param display_name: watcher display_name
        :param computed_state: watcher state
        :param output: watcher output
        """
        event = forger(
            connector="canopsis",
            connector_name="engine",
            event_type="watcher",
            source_type="component",
            component=_id,
            state=computed_state,
            output=output,
            perf_data_array=[],
            display_name=display_name)

        self.amqp_pub.canopsis_event(event)

    def sla_compute(self, watcher_id, state):
        """
        Launch the sla calcul.

        :param watcher_id: watcher id
        :param state: watcher state
        """

        # sla_tab = list(
        #     self.sla_storage.get_elements(query={'_id': watcher_id}))[0]
        # sla_tab['states'][state] = sla_tab['states'][state] + 1

        # self.sla_storage.put_element(sla_tab)

        # watcher_conf = list(
        #     self[self.WATCHER_STORAGE].get_elements(
        # query={'_id': watcher_id})
        # )[0]

        # sla = Sla(self[self.WATCHER_STORAGE],
        #           'test/de/rk/on/verra/plus/tard',
        #           watcher_conf['sla_output_tpl'],
        #           watcher_conf['sla_timewindow'],
        #           watcher_conf['sla_warning'],
        #           watcher_conf['alert_level'],
        #           watcher_conf['display_name'])

        # self.logger.critical('{0}'.format((
        #     sla_tab['states']/
        #     (sla_tab['states'][1] +
        #      sla_tab['states'][2] +
        #      sla_tab['states'][3]))))
        pass

    @staticmethod
    def worst_state(nb_crit, nb_major, nb_minor):
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

        return 0
