from test_base import BaseApiTest
from canopsis.context_graph.manager import ContextGraph


class Alarm:


     def __init__(self):
         pass


class BaseAlarm(BaseApiTest):

    URL_GET_ALARM = "/alerts/get-alarms"

    PARAMS_GET_ALARMS = {"filter": None,
                         "opened": None,
                         "resolved": None,
                         "lookups": ["pbehaviors","linklist"],
                         "sort_key": "v.state.t", # not really important
                         "sort_dir": "DESC", # not really important
                         "limit": 1000,
                         "skip": 0,
                         "natural_search": None}

    def create_alarm(self, event):
        self._send_event(event)

    def _gen_filter_from_event(self, event):
        """Generate a filter to match the entity specified by the given event.
        param event : the event
        return dict: the filter
        """
        return  {"$or" : [{"d":ContextGraph.get_id(event.__dict__)}]}

    def get_alarm_filter(self, params):
        url = self._build_url(self.URL_GET_ALARM, params)

        return self._send(url)
