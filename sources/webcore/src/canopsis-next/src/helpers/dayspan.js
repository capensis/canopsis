import moment from 'moment';
import { groupBy } from 'lodash';
import { Day, Schedule } from 'dayspan';

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
  const groupedAlarms = groupBy(alarms, alarm => moment.unix(alarm.t).startOf(groupByValue).format());

  return Object.keys(groupedAlarms).map((dateString) => {
    const dateObject = moment(dateString);
    const startDay = new Day(dateObject);
    const sum = groupedAlarms[dateString].length;

    return {
      data: {
        title: sum,
        description: filter.title,
        color: getColor(sum),
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
            hasPopover: true,
            events: groupedEvent,
          },
        },
      };
    }

    return groupedEvent[0];
  });
}
