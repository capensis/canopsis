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

from canopsis.alerts.manager import Alerts
from canopsis.alerts.selector import Selector
from canopsis.common.utils import singleton_per_scope
from canopsis.engines.core import Engine, publish
from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.old.record import Record


class engine(Engine):
    etype = 'selector'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        alert = singleton_per_scope(Alerts)
        self.storage = alert[Alerts.ALARM_STORAGE]

        self.storage_event = get_storage(
            namespace='object',
            account=Account(user="root", group="root")
        )

    def get_selectors(self):
        return self.storage_event.find({'crecord_type': 'selector'})

    def beat(self):
        with self.Lock(self, 'selector_processing') as l:
            if l.own():
                events = [selector.dump()
                          for selector in self.get_selectors()]

                for event in events:
                    self.logger.debug(u'Publish event: {0}'.format(event))
                    publish(
                        publisher=self.amqp,
                        event=event,
                        rk=self.amqp_queue,
                        exchange='amq.direct',
                        logger=self.logger
                    )

    def work(self, event, *args, **kwargs):
        self.logger.debug(u'Start method work: {0}'.format(event))

        selector = Selector(
            storage=self.storage,
            record=Record(event),
            logging_level=self.logging_level
        )

        self.logger.debug(u'Start processing selector: {0}'.format(selector.display_name))

        # Selector event have to be published when do state is true.
        if selector.dostate:
            self.logger.debug(u'Start get alert data')
            alert = selector.alert()
            selector.save(alert)
            self.publish_event(selector, alert)

    def publish_event(self, selector, event):
        event['selector_id'] = selector._id

        self.logger.info(
            u'Publish event: selector={} state={}'.format(
                selector.display_name,
                event['state']
            )
        )

        publish(publisher=self.amqp, event=event)

        self.logger.debug(u'Event sent')