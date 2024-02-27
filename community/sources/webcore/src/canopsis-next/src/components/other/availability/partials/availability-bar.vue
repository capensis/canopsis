<template>
  <v-layout class="gap-5" column>
    <availability-bar-filters
      :interval="query.interval"
      :show-type.sync="showType"
      :min-interval-date="minAvailableDate"
      @update:interval="updateInterval"
    />
    <availability-bar-chart
      :downtime="availability.downtime"
      :uptime="availability.uptime"
      :inactive-time="availability.inactive_time"
      :type="query.showType"
    />
  </v-layout>
</template>

<script>
import { AVAILABILITY_SHOW_TYPE, QUICK_RANGES } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';

import { metricsIntervalFilterMixin } from '@/mixins/widget/metrics/interval';

import AvailabilityBarChart from './availability-bar-chart.vue';
import AvailabilityBarFilters from './availability-bar-filters.vue';

export default {
  components: { AvailabilityBarFilters, AvailabilityBarChart },
  mixins: [metricsIntervalFilterMixin],
  props: {
    availability: {
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
      query: {
        interval: {
          from: QUICK_RANGES.today.start,
          to: QUICK_RANGES.today.stop,
        },
      },
    };
  },
  computed: {
    minAvailableDate() {
      return this.minDate
        ? convertDateToStartOfDayTimestampByTimezone(this.minDate, this.$system.timezone)
        : null;
    },
  },
  watch: {
    query: {
      immediate: true,
      handler() {
        this.$emit('update:query', this.getQuery());
      },
    },
  },
  methods: {
    getQuery() {
      return {
        interval: this.getIntervalQuery(),
      };
    },
  },
};
</script>
