#!/usr/bin/env python
# -*- coding: utf-8 -*-

import MySQLdb as mysql
import time

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
        self.handler = self.init.getHandler(self.logger)

    def load_config(self):
        """ Load configuration """

        import sys
        import os

        sys.path.append(os.path.expanduser('~/etc'))
        import glpi2amqp_conf

        try:
            self.config = {
                'mysql_host': glpi2amqp_conf.mysql_host,
                'mysql_user': glpi2amqp_conf.mysql_user,
                'mysql_pass': glpi2amqp_conf.mysql_pass,
                'mysql_db': glpi2amqp_conf.mysql_db,
                'interval': glpi2amqp_conf.interval,
            }

        except AttributeError, err:
            self.logger.error('Can\'t load configuration: {0}'.format(err))
            return False

        return True

    def __call__(self):
        self.handler.run()

        # start AMQP
        self.logger.debug('Start AMQP...')
        self.amqp = camqp()
        self.amqp.start()

        try:
            # connect to MySQL
            self.logger.debug('Connect to MySQL database {0}@{1} {2}...',
                self.config['mysql_user'],
                self.config['mysql_host'],
                self.config['mysql_db'])

            self.sql = mysql.connect(
                self.config['mysql_host'],
                self.config['mysql_user'],
                self.config['mysql_pass'],
                self.config['mysql_db'])

            # get data
            self.get_data()

            # close MySQL connection
            self.logger.debug('Close MySQL connection...')
            self.sql.close()

        except mysql.Error, err:
            self.logger.error('MySQL error #{0}: {1}'.format(err.args[0], err.args[1]))

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

                perf_data = []

                try:
                    with self.sql:
                        cursor = self.sql.cursor(mysql.cursors.DictCursor)

                        # Get informations about tickets
                        cursor.execute('''
                            SELECT
                                status,
                                COUNT(*) AS total,
                                AVG(TIME_TO_SEC(TIMEDIFF(`closedate`, `date`))) AS avgtime
                            FROM
                                glpi_tickets
                            GROUP BY `status`
                        ''')

                        rows = cursor.fetchall()

                        for row in rows:
                            # Generate data for assigned tickets
                            if row['status'] == 'assign':
                                perf_data.append({
                                    'metric': 'n_tickets_open',
                                    'value': row['total'],
                                    'unit': None,
                                    'min': 0,
                                    'max': None,
                                    'warn': None,
                                    'crit': None,
                                    'type': 'GAUGE'
                                })

                            # Generate data for closed tickets
                            elif row['status'] == 'closed':
                                perf_data.append({
                                    'metric': 'n_tickets_closed',
                                    'value': row['total'],
                                    'unit': None,
                                    'min': 0,
                                    'max': None,
                                    'warn': None,
                                    'crit': None,
                                    'type': 'GAUGE'
                                })

                                # Average time passed on tickets
                                perf_data.append({
                                    'metric': 'tickets_time_avg',
                                    'value': int(row['avgtime']),
                                    'unit': 's',
                                    'min': 0,
                                    'max': None,
                                    'warn': None,
                                    'crit': None,
                                    'type': 'GAUGE'
                                })

                except mysql.Error, err:
                    self.logger.error('MySQL error #{0}: {1}'.format(err.args[0], err.args[1]))
                    continue

                # Send data
                try:
                    self.on_log(perf_data)

                except Exception, err:
                    self.logger.error('Impossible to send log to Canopsis: \'{0}\''.format(err))
                    continue

                time.sleep(self.config['interval'])

        except Exception, err:
            self.logger.error('Exception: \'{0}\''.format(err))

    def on_log(self, data):
        """ When data are available, push them to Canopsis (via AMQP) """

        # Create event from data
        event = cevent.forger(
            connector='glpi',
            connector_name=DAEMON_NAME,
            component='glpi',
            resource='mysql',
            timestamp=int(time.time()),
            source_type='resource',
            event_type='log',
            state=0,
            perf_data_array=data
        )

        # Publish event
        self.logger.debug('Event: {0}'.format(event))

        key = cevent.get_routingkey(event)
        self.amqp.publish(event, key, self.amqp.exchange_name_events)


if __name__ == '__main__':
    connector = Connector()

    if connector.load_config():
        connector()
