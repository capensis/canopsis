from __future__ import unicode_literals

import re


class HeartBeat:

    MAX_DUR_REGEXP = "[0-9]*(s|m|h)"
    MAPPINGS_KEY = "mappings"
    MAX_DUR_KEY = "maxduration"

    def __init__(self, mappings, maxduration):
        self.mappings = mappings
        self.max_duration = maxduration

    @classmethod
    def isValid(cls, heartBeat):
        """
        Check if the heartBeat given is valid.

        In order to considered valid, a heartbeat must have two fields:
        `mappings` and `maxDuration`.

        `Mappings` is a list of item.
        An `item` is an ojects with at least one key. The key and the
        associated value are both string
        `maxduration` is a string that match the follow pattern: [0-9]*(s|m|h).
        """
        it = 0
        try:
            mappings = heartBeat[cls.MAPPINGS_KEY]
        except KeyError:
            return False, "The `mappings` field is missing."

        try:
            max_duration = heartBeat[cls.MAX_DUR_KEY]
        except KeyError:
            return False, "The `maxduration` field is missing."

        for mapping in mappings:
            for key in mapping:
                if not isinstance(key, basestring):
                    return False, "{} must be a string.".format(key)

                if not isinstance(mapping[key], basestring):
                    return False, "The value of `{0}` of the mapping object"\
                        " at index {1} must be a string.".format(key, it)

                it += 1

        if re.match(cls.MAX_DUR_REGEXP, max_duration) is not None:
            return True, ""

        return False, "The maxDuration fields does not match the" \
            " regular expression {}.".format(cls.MAX_DUR_KEY)

    def to_dict(self):
        """
        Return the representation of the current instance as a dict.

        """
        return {self.MAPPINGS_KEY: self.mappings,
                self.MAX_DUR_KEY: self.max_duration}
