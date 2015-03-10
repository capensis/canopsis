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

from canopsis.engines.core import Engine
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from time import time


class engine(Engine):

    etype = 'cancel'

    """
    This engine's goal is to compute an event cancellation.
    Event cancellation can be triggered from UI and will change event
    information within database.
    """

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

    def pre_run(self):
        self.storage = get_storage(
            namespace='events', account=Account(user="root", group="root"))

    def beat(self):
        self.logger.debug('entered in cancel BEAT')

    def work(self, event, *args, **kargs):

        if event['event_type'] in ['cancel', 'uncancel']:

            devent = self.storage.get_backend().find_one(
                {'_id': event['ref_rk']}
            )
            if devent is not None:

                update = {}
                update_query = {}

                # Preparing event for cancel,
                update['output'] = devent.get('output', '')

                # Is it a cancel ?
                if event['event_type'] == 'cancel':
                    ack_info = devent.get('ack', {})
                    # Saving status, in case cancel is undone
                    # If cancel is not in ok, it's not an alert cancellation
                    update['cancel'] = {
                        'timestamp': time(),
                        'author': event.get('author', 'unknown'),
                        'comment': event['output'],
                        'previous_status': devent.get('status', 1),
                    }
                    self.logger.info("set cancel to the event")

                    update['ack'] = ack_info
                    update['ack']['isAck'] = False
                    update['ack']['isCancel'] = True
                    # Set alert to cancelled status
                    update['status'] = 4
                    event['status'] = 4

                # Undo cancel ?
                elif event['event_type'] == 'uncancel':

                    # If event has been previously cancelled
                    if 'ack' in devent:
                        self.logger.warning(' + reseting ack')

                        update['ack'] = devent.get('ack', {})
                        update['ack']['isAck'] = True
                        update['ack']['isCancel'] = False

                        # Restore previous status
                        if 'cancel' in devent:
                            update['status'] = devent['cancel'].get(
                                'previous_status',
                                0
                            )
                            update_query['$unset'] = {'cancel': ''}

                update_query['$set'] = update

                # Update database with cancel informations
                if update:
                    self.storage.get_backend().update(
                        {'_id': event['ref_rk']},
                        update_query
                    )
