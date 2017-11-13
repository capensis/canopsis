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

import arrow
import copy
import re

from datetime import datetime as dt

import dateutil.rrule as rrulemod
from dateutil.rrule import rrule as rruleobj

from canopsis.common.enumerations import FastEnum

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


class PeriodicBehavior:
    '''PBehavior representation based on Pole Emploi specific files'''

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


class DaysOfWeek(FastEnum):
    """
    Number associated with days MUST be starting from 0 for Monday and
    increment by 1 for each next day, or calculation of Activity overlap
    will not work.

    If you need another enum like this, do NOT use this one for your projects
    and create your own.
    """

    Monday = 0
    Tuesday = 1
    Wednesday = 2
    Thursday = 3
    Friday = 4
    Saturday = 5
    Sunday = 6


class TimeUnits(FastEnum):

    Second = 1
    Minute = 60
    Hour = 60 * 60
    Day = 24 * 60 * 60


class Activity(object):

    VALID_FROM = 'valid_from'
    VALID_UNTIL = 'valid_until'
    DAYS_OF_WEEK = 'day_of_week'
    START_TIME_OF_DAY = 'start_time_of_day'
    STOP_TIME_OF_WEEK = 'stop_after_time'
    ENTITY_FILTER = 'entity_filter'
    PBEHAVIOR_IDS = 'pbehavior_ids'
    AGGREGATE_NAME = 'aggregate_name'

    def __init__(
        self, entity_filter, day_of_week, start_time_of_day,
        stop_after_time, aggregate_name=None, valid_from=None,
        valid_until=None, pbehavior_ids=None
    ):
        """
        :param dict entity_filter: mongo filter
        :param ActivityAggregate aggregate: parent activity aggregate
        :param int valid_from: timestamp
        :param int valid_until: timestamp
        :param int day_of_week: 0: Monday to 6: Sunday
        :param int start_time_of_day: number of seconds after D-Day 00:00
        :param int stop_after_time: number of seconds after D-Day 00:00
        """

        self.aggregate_name = aggregate_name

        self.valid_from = valid_from
        self.valid_until = valid_until
        self.day_of_week = day_of_week

        # seconds starting from 00:00am
        self.start_time_of_day = start_time_of_day
        # seconds starting from 00:00am
        self.stop_after_time = stop_after_time
        # mongo filter
        self.entity_filter = entity_filter
        self.pbehavior_ids = [] if pbehavior_ids is None else pbehavior_ids

    @property
    def day_of_week(self):
        """
        :rtype: set
        """
        return self._day_of_week

    @day_of_week.setter
    def day_of_week(self, day):
        if not isinstance(day, int) or day < 0 or day > 6:
            raise ValueError(
                'need an int between 0 and 6 inclusive')

        self._day_of_week = day

    @property
    def start_time_of_day(self):
        return self._start_time_of_day

    @start_time_of_day.setter
    def start_time_of_day(self, value):
        if value > 24 * TimeUnits.Hour:
            raise ValueError('value > 24 hours: unsupported')

        self._start_time_of_day = value

    @property
    def stop_after_time(self):
        return self._stop_after_time

    @stop_after_time.setter
    def stop_after_time(self, value):
        if value <= self.start_time_of_day:
            raise ValueError('stop is behind start')

        if value > 7 * TimeUnits.Day + self.start_time_of_day:
            raise ValueError('value > 7 days + {}s: unsupported'.format(
                self.start_time_of_day))

        self._stop_after_time = value

    def overlap(self, activity):
        """
        :type activity: Activity
        """
        tud = TimeUnits.Day

        # my start/Stop time
        m_st = self.day_of_week * tud + self.start_time_of_day
        m_St = (self.day_of_week * tud + self.stop_after_time) \
            % (7 * TimeUnits.Day)

        # foreign start/Stop time
        f_st = activity.day_of_week * tud + activity.start_time_of_day
        f_St = (activity.day_of_week * tud + activity.stop_after_time) \
            % (7 * TimeUnits.Day)

        start_inside = f_st >= m_st and f_st <= m_St
        stop_inside = f_St <= m_St and f_St >= m_st

        if stop_inside or start_inside:
            return True

        return False

    def __eq__(self, obj):
        """
        Returns equality on those attributes:

            valid_from
            valid_until
            day_of_week
            start_time_of_day
            stop_after_time
            entity_filter
        """

        return self.valid_from == obj.valid_from and \
            self.valid_until == obj.valid_until and \
            self.day_of_week == obj.day_of_week and \
            self.start_time_of_day == obj.start_time_of_day and \
            self.stop_after_time == obj.stop_after_time and \
            self.entity_filter == self.entity_filter

    def to_dict(self):
        return {
            self.VALID_FROM: self.valid_from,
            self.VALID_UNTIL: self.valid_until,
            self.DAYS_OF_WEEK: self.day_of_week,
            self.START_TIME_OF_DAY: self.start_time_of_day,
            self.STOP_TIME_OF_WEEK: self.stop_after_time,
            self.ENTITY_FILTER: self.entity_filter,
            self.PBEHAVIOR_IDS: self.pbehavior_ids,
            self.AGGREGATE_NAME: self.aggregate_name
        }


class ActivityAggregate(object):

    def __init__(self, name):
        super(ActivityAggregate, self).__init__()
        self._activities = []
        self.name = name

    @property
    def activities(self):
        """
        :rtype: list[Activity]
        """
        return self._activities

    def add(self, activity):
        """
        :param Activity activity:
        :raises ValueError: another activity is set on a day of that week.
        :rtype: bool
        :returns: True if activity was added, False if activity is already in
        """
        for ac in self.activities:
            if ac == activity:
                return False

        for m_activity in self.activities:
            if m_activity.overlap(activity):
                raise ValueError('activity overlap')

        activity.aggregate_name = self.name

        self._activities.append(activity)
        return True

    def to_dict(self):
        return {'_id': self.name}


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

    def __init__(self, tz='Europe/Paris'):
        self._tz = tz

    def _get_monday(self):
        """
        Get current week's monday date.
        """
        today = arrow.now(tz=self._tz)
        monday = today.shift(days=-today.weekday())
        return monday

    def _dow_to_rruledays(self, day_of_week):
        """
        :param int day_of_week:
        :rtype: list[str]
        :returns: list of days of week for BYDAY RRULE instruction
        """
        rruledow = []

        for dow in day_of_week:
            rruledow.append(DAY_OF_WEEK_RRULE_DAY[dow])

        return rruledow

    def _normalize_date(self, date, start_time, stop_time):
        start_date = arrow.get(
            dt(date.year, date.month, date.day, 0, 0, 0), tzinfo=self._tz
        ).shift(seconds=start_time)

        stop_date = arrow.get(
            dt(date.year, date.month, date.day, 0, 0, 0), tzinfo=self._tz
        ).shift(seconds=stop_time)

        return start_date, stop_date

    def _merged_activities_dates(self, activity_aggreg, from_date, days):
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
        tstart = start_date.timestamp
        tstop = stop_date.timestamp
        weekday = DAY_OF_WEEK_RRULE_DAY[start_date.weekday()]
        rrules = 'FREQ=DAILY;BYDAY={}'.format(weekday)
        return PeriodicBehavior(
            'downtime', filter_, tstart, tstop, rrule=rrules)

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


class ActivityAggregateManager(object):

    ACAGG_COLLECTION = 'default_activityaggregate'

    def __init__(self, acag_collection, activity_manager):
        """
        :type activity_manager: ActivityManager
        :type acag_collection: canopsis.common.collection.MongoCollection
        """
        self._coll = acag_collection
        self._activity_manager = activity_manager

    def store(self, aggregate):
        """
        Store an aggregate and attached activities.

        :type aggregate: ActivityAggregate
        """
        if self._coll.insert(aggregate.to_dict()) == aggregate.name:
            self._activity_manager.store(aggregate.activities)


class ActivityManager(object):
    """
    Store/get activities in/from database. Aggregates are never stored,
    they are only used to add a field in the activity so you can query
    activities grouped by aggregate.
    """

    ACTIVITY_COLLECTION = 'default_activity'

    def __init__(self, activity_collection):
        """
        :param activity_collection: MongoCollection
        :type activity_collection: canopsis.common.collection.MongoCollection
        """
        self._coll = activity_collection

    def store(self, activities):
        """
        :type activities: list[Activity]
        """
        activities = [ac.to_dict() for ac in activities]
        return self._coll.insert(activities)

    def del_by_aggregate_name(self, aggregate_name):
        """
        :type aggregate_name: str
        """
        return self._coll.remove({'aggregate_name': aggregate_name})

    def get(self, _id):
        """
        :rtype: Activity
        """
        act = self._coll.find_one({'_id': _id})
        act.pop('_id')

        return Activity(**act)

    def get_by_aggregate_name(self, aggregate_name):
        """
        :param str aggregate_name:
        :rtype: list[Activity]
        """
        activities = []
        res = self._coll.find({'aggregate_name': aggregate_name})

        for r in list(res):
            r.pop('_id')
            activities.append(Activity(**r))

        return activities
