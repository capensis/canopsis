# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from copy import copy

from canopsis.common.enumerations import FastEnum


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

    AGGREGATE_NAME = 'aggregate_name'
    VALID_FROM = 'valid_from'
    VALID_UNTIL = 'valid_until'
    DAYS_OF_WEEK = 'day_of_week'
    START_TIME_OF_DAY = 'start_time_of_day'
    STOP_TIME_OF_WEEK = 'stop_after_time'
    ENTITY_FILTER = 'entity_filter'
    DBID = 'dbid'

    def __init__(
        self, entity_filter, day_of_week, start_time_of_day,
        stop_after_time, aggregate_name=None, valid_from=None,
        valid_until=None, dbid=None
    ):
        """
        :param dict entity_filter: mongo filter
        :param ActivityAggregate aggregate: parent activity aggregate
        :param int valid_from: timestamp
        :param int valid_until: timestamp
        :param int day_of_week: 0: Monday to 6: Sunday
        :param int start_time_of_day: number of seconds after D-Day 00:00
        :param int stop_after_time: number of seconds after D-Day 00:00
        :param string aggregate_name: aggregate name is the aggregate id
        :raises ValueError: initialisation failed
        """
        self.aggregate_name = aggregate_name

        self.valid_from = valid_from
        self.valid_until = valid_until
        self.day_of_week = day_of_week

        self.start_time_of_day = start_time_of_day
        self.stop_after_time = stop_after_time

        self.entity_filter = entity_filter
        self.dbid = dbid

    @staticmethod
    def __copy__(remote):
        return Activity(
            remote.entity_filter,
            remote.day_of_week,
            remote.start_time_of_day,
            remote.stop_after_time,
            aggregate_name=remote.aggregate_name,
            valid_from=remote.valid_from,
            valid_until=remote.valid_until,
            dbid=remote.dbid
        )

    @property
    def day_of_week(self):
        """
        :rtype: set
        """
        return self._day_of_week

    @day_of_week.setter
    def day_of_week(self, day):
        if not isinstance(day, int) or day < 0 or day > 6:
            raise ValueError('need an int between 0 and 6 inclusive')

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

        return (self.valid_from == obj.valid_from
                and self.valid_until == obj.valid_until
                and self.day_of_week == obj.day_of_week
                and self.start_time_of_day == obj.start_time_of_day
                and self.stop_after_time == obj.stop_after_time
                and self.entity_filter == self.entity_filter)

    def to_dict(self):
        return {
            self.VALID_FROM: self.valid_from,
            self.VALID_UNTIL: self.valid_until,
            self.DAYS_OF_WEEK: self.day_of_week,
            self.START_TIME_OF_DAY: self.start_time_of_day,
            self.STOP_TIME_OF_WEEK: self.stop_after_time,
            self.ENTITY_FILTER: self.entity_filter,
            self.AGGREGATE_NAME: self.aggregate_name
        }


class ActivityAggregate(object):

    def __init__(self, name, entity_filter, pb_ids=None):
        super(ActivityAggregate, self).__init__()
        self._activities = []
        self.name = name
        self.entity_filter = entity_filter
        self.pb_ids = [] if pb_ids is None else list(set(pb_ids))

    @property
    def activities(self):
        """
        :rtype: list[Activity]
        """
        return self._activities

    def add(self, activity):
        """
        Replace the activity (copy of) filter by the aggregate filter
        then add the activity to the aggregate.

        :param Activity activity: sic
        :raises ValueError: another activity is set on a day of that week.
        :rtype: bool
        :returns: True if activity was added, False if activity is already in
        """
        activity = copy(activity)
        activity.entity_filter = self.entity_filter
        activity.aggregate_name = self.name

        for ac in self.activities:
            if ac == activity:
                return False

        for m_activity in self.activities:
            if m_activity.overlap(activity):
                raise ValueError('activity overlap')

        self._activities.append(activity)
        return True

    def to_dict(self):
        return {
            '_id': self.name,
            'pb_ids': self.pb_ids,
            'entity_filter': self.entity_filter
        }
