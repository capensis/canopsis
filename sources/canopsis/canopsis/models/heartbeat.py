from __future__ import unicode_literals

import re

class HeartBeat:

    MAX_DUR_REGEXP = "[0-9]*(s|m|h)"
    MAPPING_KEY = "mapping"
    MAX_DUR_KEY = "maxduration"

    def __init__(self, mapping, maxDuration):
        self.mapping = mapping
        self.maxDuration = maxDuration

    @classmethod
    def isValid(cls, heartBeat):
        """
        Check if the heartBeat given is valid.

        In order to considered valid, a heartbeat must have two fields:
        `mappings` and `maxDuration`.

        `Mappings` is a list of item.
        An `item` is an ojects with at least one key. The key and the
        associated value are both string
        `maxDuration` is a string that match the follow pattern: [0-9]*(s|m|h).
        """
        it = 0
        mapping = heartBeat[cls.MAPPING_KEY]
        for key in mapping:
            if not isinstance(key, basestring):
                return False, "{} must be a string.".format(key)

            if not isinstance(mapping[key], basestring):
                return False, "The value associated to {0} of the element"\
                    " at index {1} must be a string.".format(key, it)

            it += 1

        if re.match(cls.MAX_DUR_REGEXP, heartBeat[cls.MAX_DUR_KEY])is not None:
            return True, ""

        return False, "The maxduration fields does not match the" \
            " regular expression {}.".format(cls.MAX_DUR_KEY)

    def to_dict(self):
        """
        Return the representation of the current instance as a dict.

        """

        return {self.MAPPING_KEY: self.mapping,
                self.MAX_DUR_KEY: self.maxDuration}
