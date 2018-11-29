# -*- coding: utf-8 -*-

from __future__ import print_function, unicode_literals

from sqlalchemy.exc import SQLAlchemyError as SQLError
from sqlalchemy.sql import text as SQLQuery
from sqlalchemy import create_engine as SQLEngine

from kombu import Connection as AMQPConnection
from kombu.pools import producers as AMQPProducers

from socket import error as SocketError

from ConfigParser import ConfigParser, Error as ConfigError
from argparse import ArgumentParser

import datetime
import decimal
import logging
import json
import sys
import datetime


class EmptyResultException(Exception):
    pass


def alchemyencoder(obj):
    """
    JSON encoder function for SQLAlchemy special classes.
    """

    if isinstance(obj, datetime.date):
        return obj.isoformat()

    elif isinstance(obj, decimal.Decimal):
        i = int(obj)
        f = float(obj)

        return i if i == f else f


class Application(object):

    """
        Application class:

        Attributes:
            argparser instance of ArgumentParser
            config instance of ConfigParser from ConfigParser
            engine create_engine from sqlalchemy
            conn connection to sql engine
    """

    def __init__(self, *args, **kwargs):
        super(Application, self).__init__(*args, **kwargs)

        self.argparser = ArgumentParser(
            description='SQL connector to Canopsis'
        )

        self.argparser.add_argument(
            '-c', '--config',
            type=str, nargs=1,
            help='Path to configuration file'
        )

        self.argparser.add_argument(
            '-l', '--loglevel',
            type=str, nargs=1,
            help='Log level (default: info)'
        )

        self.config = ConfigParser()

        self.logger = logging.getLogger('sql2canopsis')

        self.engine = None
        self.conn = None
        self.last_value = None
        self.last_value_name = 'sql_last_value'
        self.last_value_column = None
        self.last_value_init = None
        self.last_value_retention_file = None

    def sql_connect(self):
        """
            method to etablished a connection to the database

            :return: a boolean
        """
        try:
            dburi = self.config.get('database', 'url')

            self.logger.debug('database.url = %s', dburi)

            self.engine = SQLEngine(dburi)
            self.conn = self.engine.connect()

        except SQLError as err:
            self.logger.error(
                'Impossible to connect to database: %r',
                err, exc_info=1
            )

            return False

        except ConfigError as err:
            self.logger.error(
                'Impossible to read database.url from config: %r',
                err, exc_info=1
            )

            return False

        return True

    def sql_query(self):
        """
            method to execute sql query in the database
            :return: ResultProxy from SqlAlchemy, or None if
            query failed.
        """

        result = None

        try:
            dbquery = self.config.get('database', 'query')
            if self.last_value_name is not None:
                dbquery = dbquery.format(
                    **{self.last_value_name: self.last_value})
            q = SQLQuery(dbquery)

            result = self.conn.execute(q)

        except SQLError as err:
            self.logger.error(
                'Impossible to query database: %r',
                err, exc_info=1
            )

        except ConfigError as err:
            self.logger.error(
                'Impossible to read database.query from config: %r',
                err, exc_info=1
            )

        return result

    def fetch_next_batch(self, result_proxy):
        items = []
        keys = result_proxy.keys()

        rows = result_proxy.fetchmany(
            int(self.config.get('database', 'fetch_size')))

        if not rows:
            raise EmptyResultException()

        for row in rows:
            items.append(dict(zip(keys, row)))

        return json.loads(json.dumps(items, encoding=self.config.get('database', 'encoding'), default=alchemyencoder))

    def gen_events(self, items):
        """
            method to generate event
            :param items: sql query results.
            :return: events
        """
        try:
            events = []
            evtemplate = dict(self.config.items('event'))

            for item in items:
                event = {}
                perfdata = {}

                actualstate = 0
                info_state = []

                for key in evtemplate:

                    try:
                        key_name, key_type = key.split('.')
                        key_state = ''
                    except:
                        key_name, key_type, key_state = key.split('.')

                    if key_type.startswith('metric'):
                        if key_name not in perfdata:
                            perfdata[key_name] = {}

                    if key_type == 'constant':
                        event[key_name] = evtemplate[key]

                    elif key_type == 'value':
                        item_key = evtemplate[key]
                        event[key_name] = item[item_key]

                    elif key_type == 'metricvalue':
                        item_key = evtemplate[key]
                        perfdata[key_name]['value'] = item[item_key]

                    elif key_type == 'metrictype':
                        item_key = evtemplate[key]
                        perfdata[key_name]['type'] = evtemplate[key]

                    elif key_type == 'metric':
                        if key_state == 'value':
                            item_key = evtemplate[key]
                            info_state.append(
                                (key_name, 'value', item[item_key]))
                            perfdata[key_name]['value'] = item[item_key]
                        elif key_state == 'min':
                            info_state.append(
                                (key_name, 'min', evtemplate[key]))
                        elif key_state == 'maj':
                            info_state.append(
                                (key_name, 'maj', evtemplate[key]))
                        elif key_state == 'crit':
                            info_state.append(
                                (key_name, 'crit', evtemplate[key]))

                if info_state:
                    test = self.get_state(info_state, actualstate)
                    if (test != -1):
                        actualstate = test
                        event['state'] = test

                event['perf_data_array'] = []

                for key in perfdata:
                    metric = perfdata[key]
                    metric['metric'] = key

                    event['perf_data_array'].append(metric)

                events.append(event)

        except ConfigError as err:
            self.logger.error(
                'Impossible to read event section from config: %r',
                err, exc_info=1
            )

            events = []

        return events

    def states_objects(self, info_state):
        """
            methode to get differents name of things to get states

            :param info_state: a list with tuples

            :return: a list with names
        """
        tmp = []
        for infos in info_state:
            if self.not_in(infos[0], tmp):
                tmp.append(infos[0])
        return tmp

    def not_in(self, name, alist):
        """
        method to check if an element is in a list
        :param name: an element
        :param alist: a list of element
        :return: a boolean
        """
        for things in alist:
            if things == name:
                return False
        return True

    def get_state(self, info_state, actualstate):
        """
            get information to have different values to generate a state

            :param info_state: a list with informations
            :param actualstate: the state before the check

            :return: a state to put in the event
        """

        objects = self.states_objects(info_state)

        states = []

        for obj in objects:

            tmp = []
            minor_min = ''
            minor_max = ''
            maj_min = ''
            maj_max = ''
            crit_min = ''
            crit_max = ''
            value = ''

            for name in info_state:
                if obj == name[0]:
                    tmp.append(name)

            for tup in tmp:
                if tup[1] == 'min':
                    if ':' in tup[2]:
                        minor_min = tup[2].replace(':', '')
                    else:
                        minor_max = tup[2]

                elif tup[1] == 'value':
                    value = tup[2]

                elif tup[1] == 'maj':
                    if ':' in tup[2]:
                        maj_min = tup[2].replace(':', '')
                    else:
                        maj_max = tup[2]

                elif tup[1] == 'crit':
                    if ':' in tup[2]:
                        crit_min = tup[2].replace(':', '')
                    else:
                        crit_max = tup[2]

            states.append(self.state(
                value,
                minor_min,
                minor_max,
                maj_min,
                maj_max,
                crit_min,
                crit_max))

        if states:
            return max(states)
        else:
            return -1

    def state(
            self,
            value,
            minor_min,
            minor_max,
            maj_min,
            maj_max,
            crit_min,
            crit_max):
        """
        generate a state with informations

        :param value:
        :param minor_min
        :param minor_max
        :param maj_min
        :param maj_max
        :param crit_min
        :param crit_max

        :return: a state 0, 1, 2, 3
        """

        tmp = []

        if (not (minor_min == '')):
            if float(value) < float(minor_min):
                if (maj_min != ''):
                    if float(value) < float(maj_min):
                        if (crit_min != ''):
                            if float(value) < float(crit_min):
                                tmp.append(3)
                            else:
                                tmp.append(2)
                        else:
                            tmp.append(2)
                    else:
                        tmp.append(1)
                elif (crit_min != ''):
                    if float(value) < float(crit_min):
                        tmp.append(3)
                    else:
                        tmp.append(1)
                else:
                    tmp.append(1)
            else:
                tmp.append(0)
        elif(maj_min != ''):
            if float(value) < float(maj_min):
                if (crit_min != ''):
                    if float(value) < float(crit_min):
                        tmp.append(3)
                    else:
                        tmp.append(2)
                else:
                    tmp.append(2)
            else:
                tmp.append(0)
        elif (crit_min != ''):
            if float(value) < float(crit_min):
                tmp.append(3)
            else:
                tmp.append(0)

        if (minor_max != ''):
            if float(value) > float(minor_max):
                if (maj_max != ''):
                    if float(value) > float(maj_max):
                        if (crit_max != ''):
                            if float(value) > float(crit_max):
                                tmp.append(3)
                            else:
                                tmp.append(2)
                        else:
                            tmp.append(2)
                    else:
                        tmp.append(1)
                elif (crit_max != ''):
                    if float(value) > float(crit_max):
                        tmp.append(3)
                    else:
                        tmp.append(1)
                else:
                    tmp.append(1)
            else:
                tmp.append(0)
        elif(maj_max != ''):
            if float(value) > float(maj_max):
                if (crit_max != ''):
                    if float(value) > float(crit_max):
                        tmp.append(3)
                    else:
                        tmp.append(2)
                else:
                    tmp.append(2)
            else:
                tmp.append(0)
        elif(crit_max != ''):
            if float(value) > float(crit_max):
                tmp.append(3)
            else:
                tmp.append(0)

        if tmp:
            return max(tmp)
        else:
            return -1

    def amqp_publish(self, events):
        """
        publish an event on amqp

        :param events: an event to publish
        :return: a boolean
        """
        try:
            amqpuri = self.config.get('amqp', 'url')

            self.logger.debug('amqp.url = %s', amqpuri)

            with AMQPConnection(amqpuri) as conn:
                with AMQPProducers[conn].acquire(block=True) as producer:
                    for event in events:
                        self.logger.debug(
                            'Event: %s',
                            json.dumps(event, indent=2)
                        )

                        rk = '{0}.{1}.{2}.{3}.{4}'.format(
                            event['connector'],
                            event['connector_name'],
                            event['event_type'],
                            event['source_type'],
                            event['component']
                        )

                        if event['source_type'] == 'resource':
                            rk = '{0}.{1}'.format(rk, event['resource'])

                            # ensure string
                            event['resource'] = '%s' % event['resource']

                        producer.publish(
                            event, serializer='json',
                            exchange='canopsis.events', routing_key=rk
                        )

                        self.logger.info('Event %s sent', rk)

        except ConfigError as err:
            self.logger.error(
                'Impossible to read amqp.url from config: %s',
                err, exc_info=1
            )

            return False

        except SocketError as err:
            self.logger.error(
                'Impossible to connect to AMQP server: %r',
                err, exc_info=1
            )

            return False

        except Exception as err:
            self.logger.error('Unhandled error: %r', err, exc_info=1)

            return False

        return True

    def update_last_value(self):
        """
        Updates the following values from configuration:

            self.last_value_fpath
            self.last_value_column
            self.last_value

        This function does nothin in use_last_value is false, from configuration.
        """
        if self.config.getboolean('database', 'use_last_value'):
            self.last_value_fpath = self.config.get(
                'database', 'last_value_retention_file')
            self.last_value_column = self.config.get(
                'database', 'last_value_column')

            if self.last_value_fpath is not None:
                try:
                    with open(self.last_value_fpath, 'r') as last_value_file:
                        self.last_value = last_value_file.read().rstrip()
                except IOError:
                    self.last_value = 0

    def write_last_value(self, last_row):
        """
        Writes the last value from last_row[self.last_value_column] into the retention file.

        :param last_row: an item retrieved from self.fetch_next_batch().
        """
        if self.config.getboolean('database', 'use_last_value'):
            last_result = last_row[self.last_value_column]
            with open(self.last_value_fpath, 'w') as last_value_file:
                last_value_file.write(last_result)

    def __call__(self):
        args = self.argparser.parse_args()

        if not args.config:
            self.argparser.print_help()
            return 1

        loglevel = args.loglevel[0] if args.loglevel else 'info'
        loglevel = getattr(logging, loglevel.upper())

        self.logger.setLevel(loglevel)

        try:
            self.logger.info('Reading configuration: %s', args.config)
            self.config.read(args.config)

        except ConfigError as err:
            self.logger.error(
                'Impossible to parse configuration %s: %r',
                args.config, err,
                exc_info=1
            )
            return 1

        if not self.sql_connect():
            return 1

        self.update_last_value()

        result_proxy = self.sql_query()
        last_result = None

        try:
            while True:

                results = self.fetch_next_batch(result_proxy)
                events = self.gen_events(results)

                for event in events:
                    print(json.dumps(event, indent=2))

                if not self.amqp_publish(events):
                    self.logger.error('Proplem during event publishing')

                self.write_last_value(results[-1])

        except EmptyResultException, ex:
            pass

        return 0


if __name__ == '__main__':
    app = Application()
    sys.exit(app())
