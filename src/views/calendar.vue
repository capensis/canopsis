<template lang="pug">
  stats-calendar(:calendar="calendar", :events="events", @change="change")
</template>

<script>
import moment from 'moment';
// import { rrulestr } from 'rrule';
import { Calendar, Schedule, Day } from 'dayspan';

import StatsCalendar from '@/components/other/stats/calendar/calendar.vue';

import entitiesAlarmMixin from '@/mixins/entities/alarm';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

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
    await this.fetchContextEntitiesList({
      widgetId: this.widget._id,
      params: {
        start: 0,
        limit: 50,
        _filter: {
          $and: [
            {
              $or: [
                { name: { $regex: 'component', $options: 'i' } },
                { type: { $regex: 'component', $options: 'i' } },
              ],
            },
          ],
        },
      },
    });

    const filter = {
      $or: [{
        connector_name: {
          $in: this.contextEntities.map(v => v.name),
        },
      }],
    };

    await this.fetchAlarmsList({
      widgetId: this.widget._id,
      params: {
        filter,
      },
    });

    setTimeout(() => {
      const dayObject = new Day(moment());

      this.events.push({
        data: {
          title: 'PBEHAVIOR',
          description: 'Something',
          color: '#3F51B5',
        },
        schedule: new Schedule({
          on: dayObject,
          times: [dayObject.asTime()],
        }),
      });
    }, 5000);
  },
  methods: {
    change({ calendar }) {
      this.events = [{
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
      }];
    },
  },
};
</script>
