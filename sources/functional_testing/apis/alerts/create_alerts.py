from __future__ import unicode_literals

import json
import unittest

from test_base import Method, HTTP, Event
from base_test import Alarm

CONNECTOR = "ft_connector"
CONNECTOR_NAME = "ft_connector"
COMPONENT = "ft_component"
RESOURCE = "ft_resource"

STATE_OK = 0
STATE_WARNING = 1
STATE_CRIITCAL = 2
STATE_UNKNOWN = 3

STATE_TYPE_SOFT = 0
STATE_TYPE_HARD = 1

class TestCreateAlert():

    def setUp(self):
        self._authenticate()  # default setup

    def test_create_alerts_state_ok(self):

        event = Event(connector=CONNECTOR,
                      connector_name=CONNECTOR_NAME,
                      component=COMPONENT,
                      resource=RESOURCE,
                      state=STATE_OK)

        self.create_alarm(event)

    def test_create_alerts_state_warning(self):
        pass

    def test_create_alerts_state_critical(self):
        pass

    def test_create_alerts_state_unknown(self):
        pass
