from test_base import BaseApiTest

class Alarm:


     def __init__(self):
         pass


class BaseAlarm(BaseApiTest):

    def create_alarm(event):
        self._send_event(event)

    def get_alarm(filter_):
        raise NotImplementedError()

    def delete_alarm(alarm):
        raise NotImplementedError()
