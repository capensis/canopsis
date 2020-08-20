import Vue from 'vue';
import moment from 'moment';
import { groupBy } from 'lodash';
import { Day, Schedule, Constants, Op, DaySpan } from 'dayspan';

import { COUNTER_GROUPING_TYPES } from '@/constants';

/**
 * Get Schedule instance for a span
 *
 * @param {DaySpan} span
 * @returns {Schedule}
 */
export function getScheduleForSpan(span) {
  const { start } = span;
  const minutes = span.minutes(Op.UP);
  const isDay = (minutes % Constants.MINUTES_IN_DAY) === 0;

  if (isDay && span.start.isStart()) {
    return Schedule.forDay(start, span.days(Op.UP));
  }

  const isHour = (minutes % Constants.MINUTES_IN_HOUR) === 0;
  const duration = isHour ? minutes / Constants.MINUTES_IN_HOUR : minutes;
  const durationUnit = isHour ? 'hours' : 'minutes';

  return Schedule.forTime(start, start.asTime(), duration, durationUnit);
}

/**
 * Convert calendar events to grouped calendar events
 *
 * @param {Array} alarms
 * @param {string} [groupByValue='hour']
 * @param {Function} [getColor=() => {}]
 * @returns []
 */
export function convertEventsToGroupedEvents({ events, groupByValue = 'hour', getColor = () => '#fff' }) {
  const groupedEvents = groupBy(events, event => event.schedule.start.date.clone().startOf(groupByValue).format());

  return Object.keys(groupedEvents).map((dateString) => {
    const groupedEvent = groupedEvents[dateString];

    if (groupedEvent.length > 1) {
      const total = groupedEvent.reduce((acc, event) => acc + event.data.total, 0);

      return {
        ...groupedEvent[0],

        data: {
          total,

          title: total,
          color: getColor(total),
          hasPopover: true,
          events: groupedEvent,
        },
      };
    }

    return groupedEvent[0];
  });
}

/**
 * Convert counter group item to dayspan Event item
 *
 * @param {number} timestamp
 * @param {Object} counterGroup
 * @param {Object} filter
 * @param {string} [grouping = COUNTER_GROUPING_TYPES.hour]
 * @param {Function} [getColor = () => '#fff']
 * @returns {Event}
 */
export function convertCounterGroupToEvent({
  timestamp,
  counterGroup,
  filter,
  grouping = COUNTER_GROUPING_TYPES.hour,
  getColor = () => '#fff',
}) {
  const { total } = counterGroup;
  const startMoment = moment.unix(Number(timestamp));
  const endMoment = startMoment.clone().endOf(grouping);
  const startDay = new Day(startMoment);
  const endDay = new Day(endMoment);
  const daySpan = new DaySpan(startDay, endDay);
  const schedule = getScheduleForSpan(daySpan);

  schedule.adjustDefinedSpan(true);

  return {
    schedule,

    data: {
      ...Vue.$dayspan.getDefaultEventDetails(),

      color: getColor(total),
      title: total,
      description: filter.title,
      filter,
      total,
    },
  };
}
