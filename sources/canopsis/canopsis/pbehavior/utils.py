"""
Utils for pbehaviors.
"""

from dateutil.rrule import rrulestr
from time import time

def check_valid_rrule(rrule, tstart):
    """
    Check for RRULE validity.

    Raises error when invalid with a more explicit message.
    :param str rrule: rrule as string
    :param int tstart: tstart timestamp
    :raises ValueError: rrule is invalid
    :returns: True if ok, raises exception if invalid
    :rtype: bool
    """
    if rrule == '' or rrule is None:
        return True
    
    if tstart is None:
        tstart = time.now()

    try:
        rrulestr(rrule, dtstart=tstart)
    except ValueError as ex:
        raise ValueError('Invalid RRULE: {}'.format(ex))

    return True
