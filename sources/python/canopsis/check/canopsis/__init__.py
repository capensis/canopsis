from canopsis.event import Event


class Check(Event):
    """
    Manage checking event with state and status information
    """

    STATE = 'state'  #: event state field name
    STATUS = 'status'  #: event status field name
    EVENT_TYPE = 'check'  #: check event type

    OK = 0  #: ok state value
    MINOR = 1  #: minor state value
    MAJOR = 2  #: major state value
    CRITICAL = 3  #: critical state value

    def __init__(self, source, state, status, meta):

        super(Event, self).__init__(
            source=source,
            data={
                Check.STATE: state,
                Check.STATUS: status
            },
            meta=meta
        )

    @property
    def state(self):
        return self.data[Check.STATE]

    @state.setter
    def state(self, value):
        self.data[Check.STATE] = value

    @property
    def status(self):
        return self.data[Check.STATUS]

    @status.setter
    def status(self, value):
        self.data[Check.STATUS] = value
