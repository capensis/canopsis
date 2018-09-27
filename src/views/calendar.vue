<template lang="pug">
  stats-calendar(calendar.sync="calendar")
</template>

<script>
// import moment from 'moment';
// import { rrulestr } from 'rrule';
import { Calendar } from 'dayspan';

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
    };
  },
  computed: {
    calendar: {
      get() {
        return Calendar.months();
      },
      set() {
        // console.log(value);
      },
    },
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
  },
};
</script>
