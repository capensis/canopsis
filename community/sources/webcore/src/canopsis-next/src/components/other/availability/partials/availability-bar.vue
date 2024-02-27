<template>
  <v-layout class="gap-5" column>
    <availability-bar-filters
      :interval="query.interval"
      :show-type.sync="showType"
      :min-interval-date="minAvailableDate"
      @update:interval="$emit('update:interval', $event)"
    />
    <availability-bar-chart
      :downtime="availability.downtime"
      :uptime="availability.uptime"
      :inactive-time="availability.inactive_time"
      :type="showType"
    />
  </v-layout>
</template>

<script>
import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';

import AvailabilityBarChart from './availability-bar-chart.vue';
import AvailabilityBarFilters from './availability-bar-filters.vue';

export default {
  inject: ['$system'],
  components: { AvailabilityBarFilters, AvailabilityBarChart },
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
  },
  data() {
    return {
      showType: AVAILABILITY_SHOW_TYPE.percent,
    };
  },
  computed: {
    minAvailableDate() {
      return this.minDate
        ? convertDateToStartOfDayTimestampByTimezone(this.minDate, this.$system.timezone)
        : null;
    },
  },
};
</script>
