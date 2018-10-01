<template lang="pug">
  stats-calendar(:calendar="calendar", :events="events", @change="change")
</template>

<script>
import moment from 'moment';
// import { rrulestr } from 'rrule';
import { Calendar, Schedule, Day, Time } from 'dayspan';
import { createNamespacedHelpers } from 'vuex';
import randomColor from 'randomcolor';

import StatsCalendar from '@/components/other/stats/calendar/calendar.vue';

import entitiesAlarmMixin from '@/mixins/entities/alarm';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

const { mapActions: contextEntityMapActions } = createNamespacedHelpers('entity');
const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');


export default {
  components: { StatsCalendar },
  mixins: [entitiesAlarmMixin, entitiesContextEntityMixin],
  data() {
    return {
      widget: {
        _id: 'asd',
      },
      calendar: Calendar.months(),
      events: [],
    };
  },
  async mounted() {
    const { entities } = await this.fetchContextEntitiesListWithoutStore({
      params: {
        start: 0,
        limit: 50,
        _filter: {
          $and: [
            {
              $or: [
                { name: { $regex: 'engine', $options: 'i' } },
              ],
            },
          ],
        },
      },
    });

    const filter = {
      $or: [{
        connector_name: {
          $in: entities.map(v => v.name),
        },
      }],
    };

    const { alarms } = await this.fetchAlarmsListWithoutStore({
      params: {
        filter,
        skip: 10,
        limit: 15,
      },
    });

    this.events = alarms.reduce((acc, alarm) => {
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
  },
  methods: {
    ...contextEntityMapActions({
      fetchContextEntitiesListWithoutStore: 'fetchListWithoutStore',
    }),
    ...alarmMapActions({
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    change(/* { calendar } */) {
      /* this.events = [{
        data: {
          title: 'START',
          description: 'Something',
          color: '#3F51B5',
        },
        schedule: new Schedule({
          on: calendar.start,
          times: [calendar.start.asTime()],
        }),
      }, {
        data: {
          title: 'END',
          description: 'Something',
          color: '#3F51B5',
        },
        schedule: new Schedule({
          on: calendar.end,
          times: [calendar.end.asTime()],
        }),
      }]; */
    },
  },
};
</script>
