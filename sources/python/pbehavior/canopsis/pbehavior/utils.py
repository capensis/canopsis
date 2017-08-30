from dateutil.rrule import rrulestr


def check_valid_rrule(rrule):
    """
    Check for RRULE validity.

    Raises error when invalid with a more explicit message.
    :param str rrule: rrule as string
    :raises ValueError: rrule is invalid
    :returns: True|False, (None|'reason')|Exception
    :rtype bool: True if ok, raises exception if invalid
    """
    if rrule == '' or rrule == None:
        return True

    try:
        rrulestr(rrule)
    except ValueError, ex:
        raise ValueError('Invalid RRULE: {}'.format(ex))

    return True
