# -*- coding: utf-8 -*-
"""
Enregistre des temps d’activité sur des entités. Exemple :


Une entité "application" représentant une application,
est active du Lundi au Vendredi de 8h à 17h.

Cela veut donc dire que tous les jours, de 17h à 8h au jour suivant,
l’application (l’entité) est inactive.


À partir de là il est facile de générer des pbehavior
et de récupérer l’information d’activité afin de l’afficher
dans un calendrier par exemple.
"""

from __future__ import unicode_literals

import arrow
import copy
import re

from datetime import datetime as dt
import dateutil.rrule as rrulemod
from dateutil.rrule import rrule as rruleobj

DAY_OF_WEEK_RRULE_DAY = {
    0: 'MO',
    1: 'TU',
    2: 'WE',
    3: 'TH',
    4: 'FR',
    5: 'SA',
    6: 'SU'
}

RE_EXCLUDED_BYDAY = re.compile('^(.*)BYDAY=[A-Z,]+(.*)$', re.I)


class PeriodicBehavior(object):
    """PBehavior representation based on Pole Emploi specific files"""

    def __init__(
        self, name, filter_, tstart, tstop,
        tzinfo='Europe/Paris', rrule=""
    ):
        self.connector = "canopsis"
        self.name = name
        self.filter_ = filter_
        self.author = "PE import script"
        self.enabled = True  # FONCTION DE LA RRULE
        self.comments = None
        self.tzinfo = tzinfo

        self.connector_name = "canopsis"
        self.rrule = rrule
        # obligatoire. Si rrule présente, seule l'heure est utilisée,
        # pour définir la plage horaire d'application du pbehavior
        self.tstart = tstart
        # obligatoire. Si rrule présente, seule l'heure est utilisée,
        # pour définir la plage horaire d'application du pbehavior
        self.tstop = tstop

    def __str__(self):
        return '{}->{}/{}'.format(
            arrow.get(self.tstart, tzinfo=self.tzinfo),
            arrow.get(self.tstop, tzinfo=self.tzinfo),
            self.rrule
        )

    def to_dict(self):
        return {
            'name': self.name,
            'filter': self.filter_,
            'author': self.author,
            'enabled': self.enabled,
            'comments': self.comments,
            'connector_name': self.connector_name,
            'rrule': self.rrule,
            'tstart': self.tstart,
            'tstop': self.tstop,
            'tzinfo': self.tzinfo
        }


class TimeDuration(object):

    def __init__(self, hour, minute, second, duration):
        """
        :param hour int: 24h format
        :param minute int: minutes
        :param second int: seconds
        :param duration int: seconds
        """

        self.hour = hour
        self.minute = minute
        self.second = second
        self.duration = duration

    def __eq__(self, obj):
        return self.hour == obj.hour and \
            self.minute == obj.minute and \
            self.second == obj.second and \
            self.duration == obj.duration

    def __hash__(self):
        return hash(str(self))

    def __str__(self):
        return '{}:{}:{}/{}'.format(
            self.hour, self.minute, self.second, self.duration)


class PBehaviorGenerator(object):

    def __init__(self, tz='UTC'):
        self._tz = tz

    def _get_monday(self):
        """
        Get current week's monday date.
        :rtype: datetime
        """
        today = arrow.now(tz=self._tz)
        monday = today.shift(days=-today.weekday())

        return monday

    def _dow_to_rruledays(self, day_of_week):
        """
        :param int day_of_week:
        :returns: list of days of week for BYDAY RRULE instruction
        :rtype: list[str]
        """
        rruledow = []

        for dow in day_of_week:
            rruledow.append(DAY_OF_WEEK_RRULE_DAY[dow])

        return rruledow

    def _normalize_date(self, date, start_time, stop_time):
        """
        :rtype: (datetime, datetime)
        """
        start_date = arrow.get(
            dt(date.year, date.month, date.day, 0, 0, 0), tzinfo=self._tz
        ).shift(seconds=start_time)

        stop_date = arrow.get(
            dt(date.year, date.month, date.day, 0, 0, 0), tzinfo=self._tz
        ).shift(seconds=stop_time)

        return start_date, stop_date

    def _merged_activities_dates(self, activity_aggreg, from_date, days):
        """
        :rtype: [(datetime, datetime)]
        """
        activities = activity_aggreg.activities
        dts = from_date.shift(days=-1).datetime
        dtS = from_date.shift(days=days - 1).datetime
        dates = []

        for activity in activities:
            rrule = rruleobj(
                rrulemod.DAILY,
                dtstart=dts,
                byweekday=activity.day_of_week
            )

            for d in rrule.between(dts, dtS, inc=True):
                ndates = self._normalize_date(
                    d, activity.start_time_of_day, activity.stop_after_time,
                )

                dates.extend(ndates)

        return sorted(dates)

    def _merge_dates(self, dates):
        """
        Supprime les dates qui se suivent.

        Exemple :

        01/01/1970 - 00:00:00
        01/01/1970 - 00:01:00
        01/01/1970 - 00:01:00


        Devient :

        01/01/1970 - 00:00:00
        01/01/1970 - 00:01:00

        :param list dates: an even list of dates.
        :raises IndexError: odd list
        """
        mdates = copy.copy(dates)
        i = 0

        try:
            while i < len(mdates):
                if mdates[i] == mdates[i + 1]:
                    mdates.pop(i + 1)
                else:
                    i += 1
        except IndexError:
            pass

        return mdates

    def _merge_pbehavior(self, pbehaviors):
        tds_dow = {}
        dow_pb = {}

        for pb in pbehaviors:
            d_start = arrow.get(pb.tstart, tzinfo=self._tz)

            td_start = TimeDuration(
                d_start.hour,
                d_start.minute,
                d_start.second,
                pb.tstop - pb.tstart
            )

            start_dow = d_start.weekday()

            if td_start not in tds_dow:
                tds_dow[td_start] = [start_dow]
                dow_pb[td_start] = pb
            else:
                tds_dow[td_start].append(start_dow)

        for td, dow in tds_dow.items():
            pb = dow_pb[td]
            byday = ','.join(self._dow_to_rruledays(dow))
            rrule_parts = RE_EXCLUDED_BYDAY.match(pb.rrule)

            if rrule_parts:
                pb.rrule = '{}BYDAY={}{}'.format(
                    rrule_parts.group(1), byday, rrule_parts.group(2)
                )

        return dow_pb.values()

    def _generate_pbehavior(self, filter_, start_date, stop_date):
        """
        :rtype: PeriodicBehavior
        """
        weekday = DAY_OF_WEEK_RRULE_DAY[start_date.weekday()]
        rrules = 'FREQ=DAILY;BYDAY={}'.format(weekday)

        return PeriodicBehavior(
            name='downtime',
            filter_=filter_,
            tstart=start_date.timestamp,
            tstop=stop_date.timestamp,
            rrule=rrules
        )

    def activities_to_pbehaviors(self, activities, start_date, days=7):
        """
        :rtype: list[PeriodicBehavior]
        """
        active_dates = self._merged_activities_dates(
            activities, start_date, days=days)

        inactive_dates = active_dates[1:]
        inactive_dates.append(active_dates[0].shift(days=days))

        m_inactive_dates = self._merge_dates(inactive_dates)

        pbehaviors = []

        _l_mid = len(m_inactive_dates)
        for i in range(0, _l_mid - _l_mid % 2, 2):
            pb_start_date = m_inactive_dates[i]
            pb_stop_date = m_inactive_dates[i + 1]
            entity_filter = activities.activities[0].entity_filter
            pb = self._generate_pbehavior(
                entity_filter, pb_start_date, pb_stop_date
            )
            pbehaviors.append(pb)

        return self._merge_pbehavior(pbehaviors)
