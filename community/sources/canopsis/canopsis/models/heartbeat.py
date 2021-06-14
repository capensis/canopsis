from __future__ import unicode_literals

import re
import itertools
import uuid


class HeartBeat(object):
    """
    Heartbeat model abstraction.

    public attrs::
        `dict` pattern: a Heartbeat pattern
        `str` expected_interval: an event expected interval.
        `str` id (computable property): an ID of a Heartbeat
        derived from a pattern.
    """

    __EXPECTED_INTERVAL_REGEXP = re.compile(r"^[0-9]*(s|m|h)$")

    ID_KEY = "_id"
    PATTERN_KEY = "pattern"
    EXPECTED_INTERVAL_KEY = "expected_interval"
    NAME_KEY = "name"
    DESCRIPTION_KEY = "description"
    AUTHOR_KEY = "author"
    OUTPUT_KEY = "output"
    CREATED_KEY = "created"
    UPDATED_KEY = "updated"

    _FIELDS = (ID_KEY, PATTERN_KEY, EXPECTED_INTERVAL_KEY, NAME_KEY,
               DESCRIPTION_KEY, AUTHOR_KEY, OUTPUT_KEY,
               CREATED_KEY, UPDATED_KEY)

    def __init__(self, heartbeat_json):
        """

        :param `dict` heartbeat_json: a Heartbeat as a dict.

        :raises: (`ValueError`, ).
        """
        if not self.is_valid_heartbeat(heartbeat_json):
            raise ValueError('invalid heartbeat format')
        # self.pattern = heartbeat_json[self.PATTERN_KEY]
        # self.expected_interval = heartbeat_json[self.EXPECTED_INTERVAL_KEY]
        for key, value in heartbeat_json.items():
            if key in self._FIELDS:
                setattr(self, key, value)

    @property
    def id(self):
        """
        A Heartbeat ID.

        :returns: a Heartbeat pattern hash.
        :rtype: `str`.
        """
        return self.get_pattern_hash(self.pattern)

    def to_dict(self):
        """
        Dump Heartbeat model as dictionary.

        :rtype: `dict`.
        """
        return {
            "_id": self.id,
            self.PATTERN_KEY: self.pattern,
            self.EXPECTED_INTERVAL_KEY: self.expected_interval,
            self.NAME_KEY: self.name,
            self.DESCRIPTION_KEY: self.description,
            self.AUTHOR_KEY: self.author,
            self.OUTPUT_KEY: self.output
        }

    @staticmethod
    def get_pattern_hash(pattern):
        """
        Hash the heartbeat pattern.

        :param `dict` pattern: heartbeat pattern.
        :returns: heartbeat pattern hash.
        :rtype: `str`.
        """
        return str(uuid.uuid4())

    @staticmethod
    def validate_heartbeat_pattern(pattern):
        """
        Check if ``pattern`` is a non-empty dict and
        the all of them keys and values are strings.

        :param `dict` pattern: a Heartbeat pattern.
        :rtype: `bool`.
        """
        return bool(pattern) and isinstance(pattern, dict) and all((
            isinstance(i, basestring) for i
            in itertools.chain(*pattern.items())
        ))

    @classmethod
    def validate_expected_interval(cls, expected_interval):
        """
        Check if the expected event interval is valid.

        :param expected_interval: a string that represent the time to wait
        before an alarm is created in the case no event link to an entity is
        received.
        :return:
        """
        return isinstance(expected_interval, basestring) and \
            bool(cls.__EXPECTED_INTERVAL_REGEXP.match(expected_interval))

    @classmethod
    def is_valid_heartbeat(cls, heartbeat_json):
        """
        Check if the heartBeat given is valid.

        In order to considered valid, the `pattern` and `expected_interval`
        attributes must be valid.

        An `pattern` is an json-object with at least one key. The key and the
        associated value are both string.

        `expected_interval` is a string that match the follow pattern:
        ^[0-9]*(s|m|h)$.
        `s` means waiting XX seconds
        `m` means waiting XX minutes
        `h` means waiting XX hours

        :param `dict` heartbeat_json: a Heartbeat json-object.
        :rtype: `bool`.
        """
        return cls.PATTERN_KEY in heartbeat_json and \
            cls.EXPECTED_INTERVAL_KEY in heartbeat_json and \
            cls.validate_heartbeat_pattern(
                heartbeat_json[cls.PATTERN_KEY]) and \
            cls.validate_expected_interval(
                heartbeat_json[cls.EXPECTED_INTERVAL_KEY])
