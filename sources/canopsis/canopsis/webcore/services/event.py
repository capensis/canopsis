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
from __future__ import unicode_literals

import json

import requests
from bottle import HTTPError, request

from canopsis.alerts.manager import Alerts
from canopsis.common.amqp import AmqpPublishError
from canopsis.common.utils import ensure_iterable
from canopsis.common.ws import route
from canopsis.event.eventslogmanager import EventsLog
from canopsis.common.utils import singleton_per_scope
from canopsis.webcore.utils import gen_json_error, HTTP_ERROR
from canopsis.userinterface.manager import is_allow_allow_change_severity_to_info


def is_valid(ws, event):
    """
    Check event validity.

    :returns: True if the event is valid, False otherwise.
    """
    event_type = event.get("event_type")
    if event_type in ['changestate', 'keepstate'] and (event.get('state', None) == 0 and
                                                       not is_allow_allow_change_severity_to_info()):
        ws.logger.error("cannot set state to info with changestate/keepstate")
        return False

    return True


def transform_event(ws, am, event):
    """
    Transform an event according its properties.

    :returns: an event with transformations. the event given in paramter is not modified.
    :raises Exception: any exception that can occured during one of the transformations.
    """
    new_event = event.copy()

    # Add role in event
    event_type = new_event.get("event_type")
    if event_type in ['ack', 'ackremove', 'cancel', 'comment', 'uncancel', 'declareticket',
                      'done', 'assocticket', 'changestate', 'keepstate', 'snooze',
                      'statusinc', 'statusdec', 'stateinc', 'statedec']:
        role = get_role(ws)
        new_event['role'] = role
        ws.logger.info(
            u'Role added to the event. event_type = {}, role = {}'.format(event_type, role))

    # Long output must be a string
    long_output = new_event.get("long_output")
    if long_output is not None:
        if not isinstance(long_output, basestring):
            new_event["long_output"] = ""
            ws.logger.warn(u'Long output field is not a string : {}. Replacing it by ""'.format(
                type(long_output)))

    # Meta-alarm children and parents
    eid = None
    if 'ref_rk' in new_event:
        eid = new_event['ref_rk']
    if not eid:
        eid = new_event['component']
        if new_event.get('resource'): 
            eid = "{}/{}".format(new_event['resource'], eid)
    alarm = am.get_last_alarm_by_connector_eid(new_event['connector'], eid)
    if isinstance(alarm, dict) and 'v' in alarm and isinstance(alarm['v'], dict):
        def alarmvalue_has_list(
            x): return x in alarm['v'] and isinstance(alarm['v'][x], list)
        if alarmvalue_has_list('children'):
            new_event['ma_children'] = list(alarm['v']['children'])
        if alarmvalue_has_list('parents'):
            new_event['ma_parents'] = list(alarm['v']['parents'])

    return new_event


def get_role(ws):
    """
    Find role of the current user

    :returns: the user role in a string, None if the role can't be found.
    """
    try:
        session = request.environ.get('beaker.session', None)
        if session is None:
            ws.logger.warning(u'get_role(): Cannot retrieve beaker.session')
            return None
    except AttributeError as ae:
        ws.logger.error(
            u'get_role(): Error while getting beaker.session')
        return None

    try:
        user = session.get('user', None)
        if user is None:
            ws.logger.warning(
                u'get_role(): Cannot retrieve user field from beaker.session')
            return None
    except AttributeError as ae:
        ws.logger.error(
            u'get_role(): Error while getting user from beaker.session')
        return None

    try:
        role = user.get('role', None)
        if role is None:
            ws.logger.warning(u'get_role(): Cannot retrieve role from user')
    except AttributeError as ae:
        ws.logger.error(
            u'get_role(): Error while getting role from user')
        return None

    return role


def send_events(ws, am, events, exchange='canopsis.events'):
    events = ensure_iterable(events)

    sent_events = []
    failed_events = []
    retry_events = []

    for event in events:
        if not is_valid(ws, event):
            ws.logger.error(
                "event {}/{} is invalid".format(event.get("resource"), event.get("component")))
            failed_events.append(event)
            continue

        try:
            transformed_event = transform_event(ws, am, event)
        except Exception as e:
            ws.logger.error('Failed to transform event : {}'.format(e))
            failed_events.append(event)
            continue

        try:
            ws.amqp_pub.canopsis_event(transformed_event, exchange)
            sent_events.append(transformed_event)

        except KeyError as exc:
            ws.logger.error('bad event: {}'.format(exc))
            failed_events.append(transformed_event)

        except AmqpPublishError as exc:
            ws.logger.error('publish error: {}'.format(exc))
            retry_events.append(transformed_event)

    return {
        'sent_events': sent_events,
        'failed_events': failed_events,
        'retry_events': retry_events
    }


def exports(ws):
    el_kwargs = {
        'el_storage': EventsLog.provide_default_basics()
    }
    manager = singleton_per_scope(EventsLog, kwargs=el_kwargs)
    am = Alerts(*Alerts.provide_default_basics())

    @ws.application.post(
        '/api/v2/event'
    )
    def send_event_post():
        try:
            events = request.json
        except ValueError as verror:
            return gen_json_error({'description':
                                   'malformed JSON : {0}'.format(verror)},
                                  HTTP_ERROR)

        if events is None:
            return gen_json_error(
                {'description': 'nothing to return'},
                HTTPError
            )

        return send_events(ws, am, events)

    @route(ws.application.post, name='event', payload=['event', 'url'])
    @route(ws.application.put, name='event', payload=['event', 'url'])
    def send_event(event, url=None):
        if ws.enable_crossdomain_send_events and url is not None:
            payload = {
                'event': json.dumps(event)
            }

            response = requests.post(url, data=payload)

            if response.status_code != 200:
                api_response = json.loads(response.text)

                return (api_response['data'], api_response['total'])

            return HTTPError(response.status_code, response.text)

        return send_events(ws, am, event)

    @route(ws.application.get,
           name='eventslog/count',
           payload=['tstart', 'tstop', 'limit', 'select']
           )
    def get_event_count_per_day(tstart, tstop, limit=100, select={}):
        """ get eventslog log count for each days in a given period
            :param tstart: timestamp of the begin period
            :param tstop: timestamp of the end period
            :param limit: limit the count number per day
            :param select: filter for eventslog collection
            :return: list in which each item contains an interval and the
            related count
            :rtype: list
        """

        results = manager.get_eventlog_count_by_period(
            tstart, tstop, limit=limit, query=select
        )

        return results
