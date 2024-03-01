<template>
  <v-layout class="gap-5" column>
    <availability-filters
      :interval="query.interval"
      :show-type.sync="showType"
      :min-interval-date="minDate"
      @update:interval="$emit('update:interval', $event)"
    />
    <availability-bar-chart
      :downtime="availability.downtime"
      :uptime="availability.uptime"
      :inactive-time="availability.inactive_time"
      :show-type="showType"
    />
  </v-layout>
</template>

<script>
import { ref } from 'vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityBarChart from './availability-bar-chart.vue';
import AvailabilityFilters from './availability-filters.vue';

export default {
  components: { AvailabilityFilters, AvailabilityBarChart },
  props: {
    availability: {
      type: Object,
      required: true,
    },
    query: {
      type: Object,
      required: true,
    },
    minDate: {
      type: Number,
      required: false,
    },
    defaultShowType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
  },
  setup(props) {
    const showType = ref(props.defaultShowType);

    return {
      showType,
    };
  },
};
</script>
