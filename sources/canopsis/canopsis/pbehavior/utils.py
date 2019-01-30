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


def parse_exdate(exdate, utc=False):
    """Extract the date, time and timezone fron the given exdate.

    Return the date and time as a datetime object and the timezone as a string.
    :param str exdate: the date and timezone to parse
    :param bool utc: to return the datetime in utc instead of the timezone
    defined in the exdate
    :raises ValueError: if the timezone extracted from the exdate string is
    invalid or if the date extracted from the exdate is invalid.
    :return datetime: a timezone aware datetime"""

    data = exdate.split(" ")

    if len(data) != 3:
        raise ValueError("The exdate does not follow the pattern "
                         "({} TIMEZONE_NAME).".format(EXDATE_DATE_FORMAT))

    date_time = datetime.strptime(data[0] + " " + data[1], EXDATE_DATE_FORMAT)

    timezone = tz.gettz(data[2])

    if timezone is None:
        ValueError("Unknown timezone : {}.".format(data[2]))

    date_time = date_time.replace(tzinfo=timezone)
    if not utc:
        return date_time

    return date_time.astimezone(tz.gettz("UTC"))
