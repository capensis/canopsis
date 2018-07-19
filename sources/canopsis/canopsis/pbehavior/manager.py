# -*- coding: utf-8 -*-

from datetime import datetime
from dateutil.rrule import rrulestr

import pytz

class PBehavior(object):

    TSTART = 'tstart'
    TSTOP = 'tstop'

class PBehaviorManager(object):
    @staticmethod
    def check_active_pbehavior(timestamp, pbehavior):
        """
        For a given pbehavior (as dict) check that the given timestamp is active
        using either:

        the rrule, if any, from the pbehavior + tstart and tstop to define
        start and stop times.

        tstart and tstop alone if no rrule is available.

        All dates and times are computed using UTC timezone, so check that your
        timestamp was exported using UTC.

        :param int timestamp: timestamp to check
        :param dict pbehavior: the pbehavior
        :rtype bool:
        :returns: pbehavior is currently active or not
        """
        fromts = datetime.utcfromtimestamp
        tstart = pbehavior[PBehavior.TSTART]
        tstop = pbehavior[PBehavior.TSTOP]

        if not isinstance(tstart, (int, float)):
            return False
        if not isinstance(tstop, (int, float)):
            return False

        tz = pytz.UTC
        dtts = fromts(timestamp).replace(tzinfo=tz)
        dttstart = fromts(tstart).replace(tzinfo=tz)
        dttstop = fromts(tstop).replace(tzinfo=tz)

        dt_list = [dttstart, dttstop]
        rrule = pbehavior['rrule']
        if rrule:
            # compute the minimal date from which to start generating
            # dates from the rrule.
            # a complementary date (missing_date) is computed and added
            # at index 0 of the generated dt_list to ensure we manage
            # dates at boundaries.
            dt_tstart_date = dtts.date()
            dt_tstart_time = dttstart.time().replace(tzinfo=tz)
            dt_dtstart = datetime.combine(dt_tstart_date, dt_tstart_time)

            # dates in dt_list at 0 and 1 indexes can be equal, so we generate
            # three dates to ensure [1] - [2] will give a non-zero timedelta
            # object.
            dt_list = list(
                rrulestr(rrule, dtstart=dt_dtstart).xafter(
                    dttstart, count=3, inc=True
                )
            )


            # compute the "missing dates": dates before the rrule started to
            # generate dates so we can check for a pbehavior in the past.
            multiply = 1
            while True:
                missing_date = dt_list[0] - multiply * (dt_list[-1] - dt_list[-2])
                dt_list.insert(0, missing_date)

                if missing_date < dtts:
                    break

                multiply += 1

            delta = dttstop - dttstart

            for dt in sorted(dt_list):
                if dtts >= dt and dtts <= dt + delta:
                    return True

            return False

        else:
            if dtts >= dttstart and dtts <= dttstop:
                return True
            return False

        return False
