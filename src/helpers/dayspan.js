import moment from 'moment';
import groupBy from 'lodash/groupBy';
import { Day, Schedule } from 'dayspan';

import { STATS_CALENDAR_COLORS } from '@/constants';

/**
 * Get calendar event color by alarms count
 *
 * @param count
 * @returns {string}
 */
export function getEventColor(count) {
  if (count > 50) {
    return STATS_CALENDAR_COLORS.alarm.large;
  }

  if (count > 30) {
    return STATS_CALENDAR_COLORS.alarm.medium;
  }

  return STATS_CALENDAR_COLORS.alarm.small;
}

/**
 * Convert alarms to calendar events
 *
 * @param alarms
 * @param groupByValue
 * @param [filter={}]
 * @returns []
 */
export function convertAlarmsToEvents({
  alarms,
  groupByValue,
  filter = {},
}) {
  const groupedAlarms = groupBy(alarms, alarm => moment.unix(alarm.t).startOf(groupByValue).format());

  return Object.keys(groupedAlarms).map((dateString) => {
    const dateObject = moment(dateString);
    const startDay = new Day(dateObject);
    const sum = groupedAlarms[dateString].length;

    return {
      data: {
        title: sum,
        description: filter.title,
        color: getEventColor(sum),
        meta: {
          sum,
          filter,
          tstart: dateObject.unix(),
          tstop: dateObject.clone().endOf(groupByValue).unix(),
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
 * @param alarms
 * @param [groupByValue='hour']
 * @returns []
 */
export function convertEventsToGroupedEvents({ events, groupByValue = 'hour' }) {
  const groupedEvents = groupBy(events, event => event.schedule.start.date.clone().startOf(groupByValue).format());

  return Object.keys(groupedEvents).map((dateString) => {
    const groupedEvent = groupedEvents[dateString];

    if (groupedEvent.length > 1) {
      const sum = groupedEvent.reduce((acc, event) => acc + event.data.meta.sum, 0);

      return {
        ...groupedEvent[0],

        data: {
          title: sum,
          color: getEventColor(sum),
          meta: {
            sum,
            hasPopover: true,
            events: groupedEvent,
          },
        },
      };
    }

    return groupedEvent[0];
  });
}

export default {
  getEventColor,
  convertAlarmsToEvents,
  convertEventsToGroupedEvents,
};
