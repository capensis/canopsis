from __future__ import unicode_literals

import re
import itertools

from hashlib import md5


class HeartBeat(object):

    __EXPECTED_INTERVAL_REGEXP = re.compile(r"^[0-9]*(s|m|h)$")

    PATTERN_KEY = "pattern"
    EXPECTED_INTERVAL_KEY = "expected_interval"

    pattern = None  # type: dict
    expected_interval = None  # type: str
    id = None  # type: str

    def __init__(self, heartbeat_json):
        """

        :param `dict` heartbeat_json: a Heartbeat as a dict.

        :raises: (`ValueError`, ).
        """
        if not self.is_valid_heartbeat(heartbeat_json):
            raise ValueError('invalid heartbeat format')
        self.pattern = heartbeat_json[self.PATTERN_KEY]
        self.id = self.get_pattern_hash(self.pattern)
        self.expected_interval = heartbeat_json[self.EXPECTED_INTERVAL_KEY]

    def to_dict(self):
        """
        Dump Heartbeat model as dictionary.

        :rtype: `dict`.
        """
        return {
            "_id": self.id,
            self.PATTERN_KEY: self.pattern,
            self.EXPECTED_INTERVAL_KEY: self.expected_interval
        }

    @staticmethod
    def get_pattern_hash(pattern):
        """
        Hash the heartbeat pattern.

        :param `dict` pattern: heartbeat pattern.
        :returns: heartbeat pattern hash.
        :rtype: `str`.
        """
        checksum = md5()
        for chunk in itertools.chain(*((k, pattern[k])
                                       for k in sorted(pattern))):
            checksum.update(chunk)
        return checksum.hexdigest()

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

        `expected_interval` is a string that match the follow pattern: ^[0-9]*(s|m|h)$.
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
