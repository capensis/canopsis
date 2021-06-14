# -*- coding: utf-8 -*-

from __future__ import print_function, unicode_literals

from kombu import Connection as AMQPConnection
from kombu.pools import producers as AMQPProducers

from socket import error as SocketError
from argparse import ArgumentParser

import json
import sys
import re

if sys.version_info < (3, 0):
    from ConfigParser import ConfigParser, Error as ConfigError
else:
    from configparser import ConfigParser, Error as ConfigError


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
            description='SQL connector to Canopsis',
            epilog="Example : python sendevent2canopsisamqp.py -c sendevent2canopsisamqp.ini -p '{}'"
        )

        self.argparser.add_argument(
            '-c', '--config',
            type=str, nargs=1,
            help='Path to configuration file',
            required=True
        )

        self.argparser.add_argument(
            '-p', '--params_json',
            type=str, nargs=1,
            default=["{}"],
            help='params list in json format'
        )

        self.config = ConfigParser()
        self.engine = None
        self.conn = None

    @staticmethod
    def parseInt(s):
        """
        Convert the string as an int, if applicable.
        """
        try:
            return int(s)
        except ValueError:
            return s

    def find_p_Key(self, items, item_key):
        """
            method to find item with regex key
            :param items: dictionnary
            :param item_key: key as regex
            :return: value
        """
        i = 0
        return_value = ""
        for (path, value) in dpath.util.search(items, item_key, yielded=True):
            return_value = value
            i += 1
        if i > 1:
            print(
                'key ' + item_key + ' matches many parameters!',
                file=sys.stderr
            )
            sys.exit(0)

        return return_value

    def gen_events(self, items):
        """
            method to generate event
            :param items: sql query results.
            :return: events
        """
        try:

                events = []
                evtemplate = dict(self.config.items('event'))

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
                        event[key_name] = Application.parseInt(evtemplate[key])

                    elif key_type == 'value':
                        item_key = evtemplate[key]
                        key_regex = key_name + ".regex"
                        if key_regex in evtemplate.keys():
                            myregex = r"" + evtemplate[key_regex]
                            match = re.match(myregex, items[item_key])
                            if match:
                                event[key_name] = match.group(1)
                            else:
                                event[key_name] = items[item_key]
                        else:
                            event[key_name] = items[item_key]

                    elif key_type == 'metricvalue':
                        item_key = evtemplate[key]
                        perfdata[key_name]['value'] = items[item_key]

                    elif key_type == 'metrictype':
                        item_key = evtemplate[key]
                        perfdata[key_name]['type'] = evtemplate[key]

                    elif key_type == 'metric':
                        if key_state == 'value':
                            item_key = evtemplate[key]
                            info_state.append(
                                (key_name, 'value', items[item_key]))
                            perfdata[key_name]['value'] = items[item_key]
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
                print(
                    'Impossible to read event section from config:', err,
                    file=sys.stderr
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

            states.append(self.state(value,
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
        :param minor_min:
        :param minor_max:
        :param maj_min:
        :param maj_max:
        :param crit_min:
        :param crit_max:

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

            with AMQPConnection(amqpuri) as conn:
                with AMQPProducers[conn].acquire(block=True) as producer:
                    for event in events:
                        rk = '{0}.{1}.{2}.{3}.{4}'.format(
                            event['connector'],
                            event['connector_name'],
                            event['event_type'],
                            event['source_type'],
                            event['component']
                        )

                        if event['source_type'] == 'resource':
                            rk = '{0}.{1}'.format(rk, event['resource'])

                        producer.publish(
                            event, serializer='json',
                            exchange='canopsis.events', routing_key=rk
                        )

                        print('Event', rk, 'sent')
                        print(json.dumps(event, indent=2))

        except ConfigError as err:
            print(
                'Impossible to read amqp.url from config:', err,
                file=sys.stderr
            )
            return False

        except SocketError as err:
            print(
                'Impossible to connect to AMQP server:', err,
                file=sys.stderr
            )
            return False

        except Exception as err:
            print('Unhandled error', type(err), ':', err, file=sys.stderr)
            return False

        return True

    def __call__(self):
        args = self.argparser.parse_args()

        if not args.config:
            self.argparser.print_help()
            return 1

        try:
            self.config.read(args.config)

        except ConfigError as err:
            print(
                'Impossible to parse configuration', args.config, ':', err,
                file=sys.stderr
            )
            return 1

        cmdl_params = json.loads(args.params_json[0])
        events = self.gen_events(cmdl_params)

        for event in events:
            print(json.dumps(event, indent=2))

        if not self.amqp_publish(events):
            return 1

        return 0


if __name__ == '__main__':
    app = Application()
    sys.exit(app())
