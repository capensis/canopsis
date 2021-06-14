"""
Utils for pbehaviors.
"""

from dateutil.rrule import rrulestr
from dateutil import tz
from datetime import datetime


EXDATE_DATE_FORMAT = "%Y/%m/%d %H:%M:%S"


def check_valid_rrule(rrule):
    """
    Check for RRULE validity.

    Raises error when invalid with a more explicit message.
    :param str rrule: rrule as string
    :raises ValueError: rrule is invalid
    :returns: True if ok, raises exception if invalid
    :rtype: bool
    """
    if rrule == '' or rrule is None:
        return True

    try:
        rrulestr(rrule)
    except ValueError as ex:
        raise ValueError('Invalid RRULE: {}'.format(ex))

    return True
