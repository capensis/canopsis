#!/usr/bin/env python
# -*- coding: utf-8 -*-

import MySQLdb as mysql
from time import time

from camqp import camqp
from cinit import cinit
import cevent

DAEMON_NAME = 'glpi2amqp'

class Connector(object):
     """ Connector between GLPI and Canopsis (via AMQP) """

    def __init__(self):
        """ Initialize connector and start AMQP socket """

        self.init = cinit()
        self.logger = self.init.getLogger(DAEMON_NAME)
        self.handler = init.getHandler(self.logger)



    def __call__(self):
        self.handler.run()

        # start AMQP
        self.logger.debug('Start AMQP...')
        self.amqp = camqp()
        self.amqp.start()

        try:
            # connect to MySQL
            self.logger.debug('Connect to MySQL database...')
            self.sql = mysql.connect('localhost', 'glpi', 'gpli', 'glpi')

            # get data
            self.get_data()

            # close MySQL connection
            self.logger.debug('Close MySQL connection...')
            self.sql.close()

        except mysql.Error, err:
            self.logger.error('MySQL error #{0}: {1}'.format(e.args[0], e.args[1]))

        # stop AMQP
        self.logger.debug('Stop AMQP...')
        self.amqp.stop()
        self.amqp.join()

    def get_data(self):
        """ While the handler is running, get data from GLPI database """

        self.logger.info("Wait for data...")

        try:
            while self.handler.status():
                # Generate data

                data = {
                    'component': 'glpi',
                    'resource': 'n_tickets_closed',
                    'timestamp': time.time(),
                    'state': 0,
                    'source_type': 'resource',
                }

                # Get number of tickets closed
                try:
                    with self.sql:
                        cursor = self.sql.cursor()
                        cursor.execute('SELECT * FROM glpi_tickets WHERE status = %s', ('closed',))

                        data['output'] = '{0}'.format(cursor.rowcount)
                        data['long_output'] = data['output']

                except mysql.Error, err:
                    self.logger.error('MySQL error #{0}: {1}'.format(err.args[0], err.args[1]))
                    continue

                # Send data
                try:
                    self.on_log(data)

                except Exception, err:
                    self.logger.error('Impossible to send log to Canopsis: \'{0}\''.format(err))
                    continue

        except Exception, err:
            self.logger.error('Exception: \'{0}\''.format(err))

    def on_log(self, data):
        """ When data are available, push them to Canopsis (via AMQP) """

        # Create event from data
        event = cevent.forger(
            connector='glpi',
            connector_name=DAEMON_NAME,
            component=data['component'],
            resource=data['resource'],
            timestamp=data['timestamp'],
            source_type=data['source_type'],
            event_type='log',
            state=data['state'],
            output=data['output'],
            long_output=data['long_output']
        )

        # Publish event
        self.logger.debug('Event: {0}'.format(event))

        key = cevent.get_routingkey(event)
        self.amqp.publish(event, key, self.amqp.exchange_name_events)


if __name__ == '__main__':
    connector = Connector()
    connector()
