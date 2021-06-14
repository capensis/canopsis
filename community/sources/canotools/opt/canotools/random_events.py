import random
import sys
import json
import cevent
import time
import datetime
from kombu import Connection
from kombu.pools import producers
from random import randint, random, choice

results = []  # couple of (timestamp, load)

now = time.time()

# algorithm parameter names
COUNT = 'count'
PERIOD = 'period'
NEW = 'new'
FLUCTUATION = 'fluctuation'
THRESHOLD = 'threshold'
MAXV = 'maxv'

# connection parameter names
HOST = 'host'
PORT = 'port'
USER = 'user'
PASSWORD = 'password'
VHOST = 'vhost'
EXCHANGE = 'exchange'

# event parameter names
COMPONENT = 'component'
RESOURCE = 'resource'
METRIC = 'metric'
CRITICAL = 'critical'
WARNING = 'warning'
UNIT = 'unit'

# algorithm parameters
count = 2000
period = 60
new = 0
fluctuation = 5 * 1000 * 1000
threshold = 3
maxv = 1 * 1000 * 1000 * 1000

# connection parameters
host = "127.0.0.1"
port = 5672
user = "guest"
password = "guest"
vhost = "canopsis"
exchange = "canopsis.events"

# event parameters
component = 'component-' + str(now)
resource = 'resource-' + str(now)
metric = 'random_event'
critical = None
warning = None
unit = None

ARGVS = {
    COUNT: (int, 'old event count'),
    PERIOD: (int, 'period duration between two events'),
    NEW: (int, 'new event count'),
    HOST: (str, 'connection host'),
    USER: (str, 'connection user id'),
    PASSWORD: (str, 'connection pwd'),
    VHOST: (str, 'connection vhost'),
    EXCHANGE: (str, 'AMQP exchange name'),
    COMPONENT: (str, 'component name'),
    MAXV: (str, 'maximal event value'),
    THRESHOLD: (int, 'algorithm threshold'),
    FLUCTUATION: (float, 'algorithm fluctuation'),
    UNIT: (str, 'metric unit'),
    CRITICAL: (str, 'metric critical threshold'),
    WARNING: (str, 'metric warning threshold')
}

for arg in sys.argv[1:]:
    name, _, value = arg.partition('=')
    t = ARGVS.get(name, None)
    if t is None:
        raise Exception("Wrong argument %s, wait for arg=value where \
            arg in '%s'" % (arg, ARGVS))
    else:
        globals()[name] = t[0](value)

start_time = now - period * count

item_output = "%s (%s): %s\n "
items_to_display = \
    [item_output % (argv, globals()[argv], ARGVS[argv][1]) for argv in ARGVS]

message = reduce(lambda x, y: x + y, items_to_display)

print message

t = 0


def random_load(lastLoad, fluctuation):
    """
    Get a random value depending of a lastLoaf value and fluctuation value.
    """
    result = 0
    if random() > 0.48:
        result = lastLoad + randint(0, int(fluctuation))
        if result > maxv:
            result = maxv
    else:
        result = lastLoad - randint(0, int(fluctuation))
        if result < 0:
            result = 0

    return result

TIMESTAMP = 'timestamp'
VALUE = 'value'


def get_results(count, timestamp, lastLoad):
    """
    Get an array of couple of (timestamp, value) with input parameters.
    Count is number of couples.
    timestamp is starting resuls value timestamp.
    lastLoad is last value.
    """
    result = []

    for index in xrange(0, count):
        random_thresold = randint(0, threshold)
        random_fluctuation = randint(0, fluctuation)

        for i in xrange(random_thresold):
            currentLoad = random_load(lastLoad, random_fluctuation)
            result.append({TIMESTAMP: timestamp, VALUE: currentLoad})
            lastLoad = currentLoad
            timestamp += period
            index += 1

    return result


def print_results(_from, to, results):
    """
    Print results in a specific interval.
    """
    message = "Interval: %s -> %s, results=%s" %\
        (datetime.datetime.fromtimestamp(_from),
            datetime.datetime.fromtimestamp(to), results)

    print message


def get_event(result):
    """
    Initialize an event with input couple of (timestamp, value)
    """
    perf_data = {
        'metric': 'random_event',
        'value': result[VALUE],
        'unit': unit,
        'type': "GAUGE",
        'min': 0,
        'max': 0,
        'warn': None,
        'crit': None}

    result = cevent.forger(
        connector='connector-random_event',
        connector_name="random_event",
        event_type="check",
        source_type="component",
        component=component,
        state=0,
        state_type=1,
        timestamp=result[TIMESTAMP],
        perf_data_array=[perf_data]
    )

    return result


def send_events(producer, results, sleep=0):
    """
    Send events.
    Input producer is used in order to send events.
    Results is an array of couple of (timestamp, value) used to fill events.
    Input sleep determinates sleeping time duration between two event sending.
    """
    for index, result in enumerate(results):
        event = get_event(result)
        routing_key = "%s.%s.%s.%s.%s" % \
            (event['connector'], event['connector_name'],
                event['event_type'], event['source_type'], event['component'])
        producer.publish(
            event,
            serializer='json',
            exchange=exchange,
            routing_key=routing_key)
        time.sleep(sleep)
        if sleep:  # print pourcentage of done sent events
            message = '%s\% done. (%s/%s)' % \
                ((index * 100.0 / len(results)), index, len(results))
            print message

# back to the future (build history)
with Connection(hostname=host, userid=user, virtual_host=vhost) as conn:
    with producers[conn].acquire(block=True) as producer:

        # get old events
        results = get_results(count, start_time, 0)
        # send old events
        send_events(producer, results)
        # print old events
        print_results(start_time, now, results)
        # get lastLoad value
        lastLoad = 0 if not results else results[-1][VALUE]
        # get new events
        results = get_results(new, now, lastLoad)
        # send new events
        send_events(producer, results, period)
        # print new events
        print_results(now, now + new * period, results)
