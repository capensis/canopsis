<template>
  <v-layout column>
    <alarm-timeline-filters />
    <alarm-timeline-days :days="days" />
  </v-layout>
</template>

<script>
import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import AlarmTimelineFilters from './alarm-timeline-filters.vue';
import AlarmTimelineDays from './alarm-timeline-days.vue';

export default {
  components: {
    AlarmTimelineFilters,
    AlarmTimelineDays,
  },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
    meta: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    days() {
      return groupAlarmSteps(this.steps);
    },
  },
};
</script>
