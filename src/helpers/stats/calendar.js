import moment from 'moment';
import { rrulestr } from 'rrule';
import { Schedule, Day, Time } from 'dayspan';

import { STATS_CALENDAR_COLORS } from '@/constants';

/**
 * Convert period to calendar events
 *
 * @param {Object} eventData
 * @param {moment.Moment} start
 * @param {moment.Moment} end
 * @returns {Array}
 */
function convertPeriodToCalendarEvents(eventData, start, end) {
  const startDay = new Day(start);
  const endDay = new Day(end);
  const result = [];

  if (start < moment()) {
    if (start.isSame(end, 'day')) {
      result.push({
        data: eventData,
        schedule: new Schedule({
          on: startDay,
          times: [startDay.asTime()],
          duration: end.diff(start, 'minutes'),
          durationUnit: 'minutes',
        }),
      });
    } else {
      result.push({
        data: eventData,
        schedule: new Schedule({
          on: startDay,
          times: [startDay.asTime()],
          duration: start.clone().endOf('day').diff(start, 'minutes'),
          durationUnit: 'minutes',
        }),
      }, {
        data: eventData,
        schedule: new Schedule({
          on: endDay,
          times: [new Time(0)],
          duration: end.diff(end.clone().startOf('day'), 'minutes'),
          durationUnit: 'minutes',
        }),
      });

      const differenceInDays = end.diff(start, 'days');

      if (differenceInDays > 1) {
        result.push({
          data: eventData,
          schedule: new Schedule({
            start: new Day(start),
            end: new Day(end.clone().subtract(1, 'day')),
          }),
        });
      }
    }
  }

  return result;
}

/**
 * Convert pbehaviors entities into calendar events
 *
 * @param {Array} pbehaviors
 * @returns {Array}
 */
export function convertPbehaviorsToCalendarEvents(pbehaviors) {
  return pbehaviors.reduce((acc, pbehavior) => {
    const eventData = {
      title: pbehavior.name,
      color: STATS_CALENDAR_COLORS.pbehavior,
    };

    if (pbehavior.rrule) {
      const rrule = rrulestr(pbehavior.rrule, {
        dtstart: new Date(pbehavior.tstart),
      });

      if (!rrule.options.until) {
        rrule.options.until = new Date(pbehavior.tstop);
      }

      const events = rrule.all().map((date) => {
        const day = new Day(moment(date));

        return {
          data: eventData,
          schedule: new Schedule({
            on: day,
            times: [day.asTime()],
          }),
        };
      });

      return acc.concat(events);
    }

    const events = convertPeriodToCalendarEvents(eventData, moment(pbehavior.tstart), moment(pbehavior.tstop));

    return acc.concat(events);
  }, []);
}

/**
 * Convert alarms entities into calendar events
 *
 * @param {Array} alarms
 * @returns {Array}
 */
export function convertAlarmsToCalendarEvents(alarms) {
  return alarms.reduce((acc, alarm) => {
    const start = moment.unix(alarm.t);
    const end = alarm.v.resolved ? moment.unix(alarm.v.resolved) : moment();

    const eventData = {
      title: alarm.d,
      color: STATS_CALENDAR_COLORS.alarm,
    };

    const events = convertPeriodToCalendarEvents(eventData, start, end);

    return acc.concat(events);
  }, []);
}

export default {
  convertAlarmsToCalendarEvents,
  convertPbehaviorsToCalendarEvents,
};
