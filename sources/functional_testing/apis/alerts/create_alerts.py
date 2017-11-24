from __future__ import unicode_literals

import json
import unittest

from test_base import Method, HTTP, Event, Response
from apis.alerts.base_alert import BaseAlarm
from canopsis.context_graph.manager import ContextGraph

CONNECTOR = "ft_connector"
CONNECTOR_NAME = "ft_connector"
COMPONENT = "ft_component"
RESOURCE = "ft_resource"

STATE_OK = 0
STATE_MINOR= 1
STATE_MAJOR = 2
STATE_CRIITCAL = 3

# state type is not use anymore, so 0 is fine
STATE_TYPE = 0

STATUS_OFF = 0
STATUS_ONGOING = 1
STATUS_STEALTHY = 2
STATUS_FLAPPING = 3
STATUS_CANCEL = 4


class TestCreateAlert(BaseAlarm):

    def setUp(self):
        self._authenticate()  # default setup

    def test_create_alerts_state_ok(self):

        event = Event(connector=CONNECTOR,
                      connector_name=CONNECTOR_NAME,
                      component=COMPONENT,
                      resource=RESOURCE,
                      state=STATE_OK,
                      state_type=STATE_TYPE,
                      source_type="resource",
                      status=STATUS_OFF)

        self.create_alarm(event)
        id_ = ContextGraph.get_id(event.__dict__)
        filter_ = self._gen_filter_from_event(event)

        get_params = self.PARAMS_GET_ALARMS.copy()
        get_params["filter"] = self._gen_filter_from_event(event)
        get_params["opened"] = False
        get_params["resolved"] = False
        get_params["natural_search"] = False

        resp = self.get_alarm_filter(get_params)
        self.assertEquals(resp.status_code, HTTP.OK.value)

        data = resp.json()
        data = data[Response.DATA.value]

        self.assertListEqual(data[0][Response.ALARMS.value], [])
        self.assertEquals(data[0][Response.TOTAL.value], 0)
        self.assertEquals(data[0][Response.LAST.value], 0)
        self.assertEquals(data[0][Response.FIRST.value], 0)

        get_params["opened"] = True

        resp = self.get_alarm_filter(get_params)
        self.assertEquals(resp.status_code, HTTP.OK.value)

        data = resp.json()
        data = data[Response.DATA.value]

        self.assertListEqual(data[0][Response.ALARMS.value], [])
        self.assertEquals(data[0][Response.TOTAL.value], 0)
        self.assertEquals(data[0][Response.LAST.value], 0)
        self.assertEquals(data[0][Response.FIRST.value], 0)

    def test_create_alerts_state_warning(self):
        event = Event(connector=CONNECTOR,
                      connector_name=CONNECTOR_NAME,
                      component=COMPONENT,
                      resource=RESOURCE,
                      state=STATE_WARNING,
                      state_type=STATE_TYPE,
                      source_type="resource",
                      status=STATUS_PENDING)

        self.create_alarm(event)
        id_ = ContextGraph.get_id(event.__dict__)
        filter_ = self._gen_filter_from_event(event)

        get_params = self.PARAMS_GET_ALARMS.copy()
        get_params["filter"] = self._gen_filter_from_event(event)
        get_params["opened"] = False
        get_params["resolved"] = False
        get_params["natural_search"] = False

        resp = self.get_alarm_filter(get_params)
        self.assertEquals(resp.status_code, HTTP.OK.value)

        data = resp.json()
        print(data)
        data = data[Response.DATA.value]

        self.assertListEqual(data[0][Response.ALARMS.value], [])
        self.assertEquals(data[0][Response.TOTAL.value], 0)
        self.assertEquals(data[0][Response.LAST.value], 0)
        self.assertEquals(data[0][Response.FIRST.value], 0)

        get_params["opened"] = True

        resp = self.get_alarm_filter(get_params)
        self.assertEquals(resp.status_code, HTTP.OK.value)

        data = resp.json()
        print(data)
        data = data[Response.DATA.value]

        self.assertListEqual(data[0][Response.ALARMS.value], [])
        self.assertEquals(data[0][Response.TOTAL.value], 0)
        self.assertEquals(data[0][Response.LAST.value], 0)
        self.assertEquals(data[0][Response.FIRST.value], 0)

    def test_create_alerts_state_critical(self):
        pass

    def test_create_alerts_state_unknown(self):
        pass
