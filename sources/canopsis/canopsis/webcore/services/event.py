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

from canopsis.common.amqp import AmqpPublishError
from canopsis.common.utils import ensure_iterable
from canopsis.common.ws import route
from canopsis.event.eventslogmanager import EventsLog
from canopsis.common.utils import singleton_per_scope
from canopsis.webcore.utils import gen_json_error, HTTP_ERROR


def is_valid(ws, event):
    """
    Check event validity.

    :returns: True if the event is valid, False otherwise.
    """
    event_type = event.get("event_type")
    if event_type in ['changestate', 'keepstate'] and event.get('state', None) == 0:
        ws.logger.error("cannot set state to info with changestate/keepstate")
        return False

    return True


def transform_event(ws, event):
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

    return new_event


def get_role(ws):
    """
    Find role of the current user

    :returns: the user role in a string, None if the role can't be found.
    """
    session = request.environ.get('beaker.session', None)
    if session is None:
        ws.logger.warning(u'get_role(): Cannot retrieve beaker.session')
        return None

    if type(session) is not dict:
        ws.logger.warning(
            u'get_role(): session doesn\'t have the correct type - {}'.format(type(session)))
        return None

    user = session.get('user', None)
    if user is None:
        ws.logger.warning(
            u'get_role(): Cannot retrieve user from beaker.session')
        return None

    if type(user) is not dict:
        ws.logger.warning(
            u'get_role(): user doesn\'t have the correct type - {}'.format(type(user)))
        return None

    role = user.get('role', None)
    if role is None:
        ws.logger.warning(u'get_role(): Cannot retrieve role from user')

    return role


def send_events(ws, events, exchange='canopsis.events'):
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
            transformed_event = transform_event(ws, event)
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

        return send_events(ws, events)

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

            else:
                return HTTPError(response.status_code, response.text)

        else:
            return send_events(ws, event)

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
