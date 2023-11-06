import { groupBy } from 'lodash';

import {
  convertDateToTimestamp,
  convertDateToStartOfUnitString,
  convertDateToEndOfUnitTimestamp,
  convertDateToDateObject,
  convertDateToMoment,
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
    const dateObject = convertDateToDateObject(dateString);
    const sum = alarmsGroup.length;

    return {
      start: dateObject,
      name: sum,
      description: filter.title,
      color: getColor(sum),
      data: {
        meta: {
          sum,
          filter,
          tstart: convertDateToTimestamp(dateObject),
          tstop: convertDateToEndOfUnitTimestamp(dateObject, groupByValue),
        },
      },
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
  const groupedEvents = groupBy(events, event => convertDateToMoment(event.start).startOf(groupByValue).format());

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
 * Get class for calendar event menu by calendar event id
 *
 * @param {string} id
 * @return {string}
 */
export const getMenuClassByCalendarEvent = id => `calendar-event-id-${id}`;
