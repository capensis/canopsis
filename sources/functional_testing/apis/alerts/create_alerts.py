from __future__ import unicode_literals

import json
import unittest

from test_base import Method, HTTP, Event
from apis.alerts.base_alert import BaseAlarm
from canopsis.context_graph.manager import ContextGraph

CONNECTOR = "ft_connector"
CONNECTOR_NAME = "ft_connector"
COMPONENT = "ft_component"
RESOURCE = "ft_resource"

STATE_OK = 0
STATE_WARNING = 1
STATE_CRIITCAL = 2
STATE_UNKNOWN = 3

# state type is not use anymore, so 0 is fine
STATE_TYPE = 0

STATUS_OK = 0
STATUS_PENDING = 1
STATUS_FURTIF = 2
STATUS_BAGOT = 3
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
                      status=STATUS_OK)

        self.create_alarm(event)
        id_ = ContextGraph.get_id(event.__dict__)
        filter_ = self._gen_filter_from_event(event)

        get_params = self.PARAMS_GET_ALARMS.copy()
        get_params["filter"] = json.dumps(self._gen_filter_from_event(event))
        get_params["opened"] = False
        get_params["resolved"] = False
        get_params["natural_search"] = False

        resp = self.get_alarm_filter(get_params)

        self.assertEquals(resp.status_code, HTTP.OK.value)

        data = resp.json()
        print(data)


    def test_create_alerts_state_warning(self):
        pass

    def test_create_alerts_state_critical(self):
        pass

    def test_create_alerts_state_unknown(self):
        pass
