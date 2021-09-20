import { get, groupBy } from 'lodash';
import { Day, Schedule, Constants, Op, DaySpan } from 'dayspan';

import {
  convertDateToMoment,
  convertDateToTimestamp,
  convertDateToStartOfUnitString,
  convertDateToEndOfUnitTimestamp,
  convertDateToMomentByTimezone,
} from '@/helpers/date/date';

/**
 * Convert alarms to calendar events
 *
 * @param {Array} alarms
 * @param {string} groupByValue
 * @param {Object} [filter={}]
 * @param {Function} [getColor=() => {}]
 * @returns []
 */
export function convertAlarmsToEvents({
  alarms,
  groupByValue,
  filter = {},
  getColor = () => '#fff',
}) {
  const groupedAlarms = groupBy(alarms, alarm => convertDateToStartOfUnitString(alarm.t, groupByValue, null));

  return Object.entries(groupedAlarms).map(([dateString, alarmsGroup]) => {
    const dateObject = convertDateToMoment(dateString);
    const startDay = new Day(dateObject);
    const sum = alarmsGroup.length;

    return {
      data: {
        title: sum,
        description: filter.title,
        color: getColor(sum),
        meta: {
          sum,
          filter,
          tstart: convertDateToTimestamp(dateObject),
          tstop: convertDateToEndOfUnitTimestamp(dateObject, groupByValue),
        },
      },
      schedule: new Schedule({
        on: startDay,
        times: [startDay.asTime()],
        duration: 1,
        durationUnit: 'hours',
      }),
    };
  });
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
      const sum = groupedEvent.reduce((acc, event) => acc + event.data.meta.sum, 0);

      return {
        ...groupedEvent[0],

        data: {
          title: sum,
          color: getColor(sum),
          meta: {
            sum,
            events: groupedEvent,
          },
        },
      };
    }

    return groupedEvent[0];
  });
}

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
 * Get DaySpan instance for timestamps with timezone conversion
 *
 * @param {number} start
 * @param {number} end
 * @param {string} timezone
 * @param {boolean} [isDate = false] - It means that start and end are startOf('day') values
 * @returns {DaySpan}
 */
export function getSpanForTimestamps({
  start,
  end,
  timezone,
}) {
  const startMoment = convertDateToMomentByTimezone(start, timezone);
  const endMoment = convertDateToMomentByTimezone(end, timezone);
  const startDay = new Day(startMoment);
  const endDay = new Day(endMoment);

  return new DaySpan(startDay, endDay);
}

/**
 * Get class for calendar event menu by calendar event id
 *
 * @param {CalendarEvent} calendarEvent
 * @return {string}
 */
export function getMenuClassByCalendarEvent(calendarEvent) {
  return `ds-calendar-event-menu_${get(calendarEvent, 'event.id', 'placeholder')}`;
}
