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
from canopsis.event.selector import Selector
from canopsis.sla import Sla
from canopsis.event import get_routingkey
from time import time


class engine(Engine):
    """
        This engine's goal is to compute an aggregated information from an
        event selection. The event selection is done thanks to a filter witch
        can include event, exclude events or select them from a cfilter. The
        worst state is then computed on the selected event set and a new event
        holding this information is produced. This computation is triggered
        each time the crecord dispatcher emit a crecord event of selector type.
    """

    etype = 'selector'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        self.selectors = []
        self.nb_beat = 0
        self.thd_warn_sec_per_evt = 1.5
        self.thd_crit_sec_per_evt = 2

    def pre_run(self):
        # Load selectors
        self.storage = get_storage(
            namespace='object', account=Account(user="root", group="root"))

    def beat(self):
        self.logger.debug('entered in selector BEAT')

    def consume_dispatcher(self, event, *args, **kargs):

        selector = self.get_ready_record(event)

        if selector:

            event_id = event['_id']

            # Loads associated class
            selector = Selector(
                storage=self.storage, record=selector,
                logging_level=self.logging_level)

            name = selector.display_name

            self.logger.debug('----------SELECTOR----------\n')

            self.logger.debug(u'Selector {} found, start processing..'.format(
                name
            ))

            update_extra_fields = {}
            # Selector event have to be published when do state is true.
            if selector.dostate:
                rk, selector_event, publish_ack = selector.event()

                # Compute previous event to know if any difference next turn
                previous_metrics = {}
                for metric in selector_event['perf_data_array']:
                    previous_metrics[metric['metric']] = metric['value']
                update_extra_fields['previous_metrics'] = previous_metrics

                do_publish_event = selector.have_to_publish(selector_event)

                if do_publish_event:
                    update_extra_fields['last_publication_date'] = \
                        selector.last_publication_date

                    self.publish_event(
                        selector,
                        rk,
                        selector_event,
                        publish_ack
                    )

                # When selector computed, sla may be asked to be computed.
                if selector.dosla:

                    self.logger.debug('----------SLA----------\n')

                    # Retrieve user ui settings
                    # This template should be always set
                    template = selector.get_sla_output_tpl()
                    # Timewindow computation duration
                    timewindow = selector.get_sla_timewindow()
                    sla_warning = selector.get_sla_warning()
                    sla_critical = selector.get_sla_critical()
                    alert_level = selector.get_alert_level()
                    display_name = selector.display_name

                    rk = get_routingkey(selector_event)

                    sla = Sla(
                        self.storage,
                        rk,
                        template,
                        timewindow,
                        sla_warning,
                        sla_critical,
                        alert_level,
                        display_name,
                        logger=self.logger
                    )
                    self.publish_sla_event(
                        sla.get_event(),
                        display_name
                    )

            else:
                self.logger.debug(u'Nothing to do with selector {}'.format(
                    name
                ))

            # Update crecords informations
            self.crecord_task_complete(event_id, update_extra_fields)

        self.nb_beat += 1
        # Set record free for dispatcher engine

    def publish_sla_event(self, event, display_name):

        publish(publisher=self.amqp, event=event)

        self.logger.debug(u'published event sla selector {}'.format(
            display_name
        ))

    def publish_event(self, selector, rk, selector_event, publish_ack):

        selector_event['selector_id'] = selector._id

        self.logger.info(
            u'Ready to publish selector {} event with state {}'.format(
                selector.display_name,
                selector_event['state']
            )
        )

        if publish_ack:
            # Define a clean ack information to the event
            now = int(time())
            selector_event['ack'] = {
                'timestamp': now,
                'rk': rk,
                'author': 'canopsis',
                'comment': 'All matched event are acknowleged',
                'isAck': True
            }
            self.logger.debug(
                'Selector event is ack because ' +
                'all matched NOK event are ack'
            )
        else:
            # Define or reset ack key for selector generated event
            selector_event['ack'] = {}
            self.logger.debug('Selector event is NOT ack')

        publish(publisher=self.amqp, event=selector_event, rk=rk)

        self.logger.debug(u'published event selector {}'.format(
            selector.display_name
        ))
