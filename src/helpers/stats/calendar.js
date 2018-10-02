import moment from 'moment';
import randomColor from 'randomcolor';
import { Schedule, Day, Time } from 'dayspan';

export function convertAlarmsToCalendarEvents(alarms) {
  return alarms.reduce((acc, alarm) => {
    const color = randomColor();
    const start = moment.unix(alarm.t);
    const end = alarm.v.resolved ? moment.unix(alarm.v.resolved) : moment();
    const startDay = new Day(start);
    const endDay = new Day(end);

    const eventData = {
      title: alarm.d,
      color,
    };

    if (start < moment()) {
      if (start.isSame(end, 'day')) {
        acc.push({
          data: eventData,
          schedule: new Schedule({
            on: startDay,
            times: [startDay.asTime()],
            duration: end.diff(start, 'minutes'),
            durationUnit: 'minutes',
          }),
        });
      } else {
        acc.push({
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
          acc.push({
            data: eventData,
            schedule: new Schedule({
              start: new Day(start),
              end: new Day(end.clone().subtract(1, 'day')),
            }),
          });
        }
      }
    }

    return acc;
  }, []);
}

export default {
  convertAlarmsToCalendarEvents,
};
