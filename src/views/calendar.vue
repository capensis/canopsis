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
        limit: 10,
      },
    });

    const events = alarms.reduce((acc, alarm) => {
      const color = randomColor();
      const start = moment.unix(alarm.t);
      const end = alarm.v.resolved ? moment.unix(alarm.v.resolved) : moment();
      const startDay = new Day(start);
      const endDay = new Day(end);

      const eventData = {
        title: alarm.d,
        color,
      };

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
            duration: start.endOf('day').diff(start, 'minutes'),
            durationUnit: 'minutes',
          }),
        }, {
          data: eventData,
          schedule: new Schedule({
            on: endDay,
            times: [new Time(0)],
            duration: end.startOf('day').diff(end, 'minutes'),
            durationUnit: 'minutes',
          }),
        });

        const difference = end.diff(start, 'days');

        if (difference > 1) {
          acc.push({
            data: eventData,
            schedule: new Schedule({
              on: new Day(start.clone().add(1, 'day')),
              duration: difference - 1,
              durationUnit: 'days',
            }),
          });
        }
      }

      return acc;
    }, []);

    this.events = events;
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
