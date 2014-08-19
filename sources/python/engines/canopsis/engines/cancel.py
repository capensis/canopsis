# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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
from canopsis.old.storage import get_storage
from canopsis.old.account import Account


class engine(Engine):

    etype = 'cancel'

    """
        This engine's goal is to compute an event cancellation.
        Event cancellation can be triggered from UI and will change event information within database.
    """

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)


    def pre_run(self):
        #load selectors
        self.storage = get_storage(
            namespace='events', account=Account(user="root", group="root"))

    def beat(self):
        self.logger.debug('entered in selector BEAT')

    def work(self, event, *args, **kargs):

        if event['event_type'] in ['cancel','uncancel']:

            devent = self.storage.get_backend().find_one({'rk': event['ref_rk']})
            if devent != None:

                update = {}

                #Saving status, in case cancel is undone
                previous_status = devent.get('status', 1)
                previous_state = devent['state']
                # Cancel value means cancel for True or undo cancel for False
                cancel = {
                    'author' : event.get('author','unknown'),
                    'isCancel' : True,
                    'comment' : event['output'],
                    'previous_status': previous_status,
                    'ack': devent.get('ack', {})
                }

                # Preparing event for cancel,
                update['output'] = devent.get('output', '')

                #Is it a cancel ?
                if event['event_type'] == 'cancel':
                    self.logger.info("set cancel to the event")
                    # If cancel is not in ok state, then it is not an alert cancellation
                    update['ack'] = cancel
                    update['status'] = 4

                # Undo cancel ?
                elif event['event_type'] == 'uncancel':

                    #If event has been previously cancelled
                    if 'ack' in devent and 'ack' in devent['ack']:
                        self.logger.warning(' + reseting ack')

                        #Restore previous status
                        update['status'] = devent['ack'].get('previous_status', 0)
                        #Restore ack as previously
                        update['ack'] = devent['ack']['ack']

                #update database with cancel informations
                if update:
                    self.storage.update({'rk': event['ref_rk']}, {'$set': update})