<template>
  <v-layout>
    <chart-loader v-if="pending" :has-data="!!availability" />
    <availability-bar
      v-if="availability"
      :query="query"
      :availability="availability"
      :default-show-type="defaultShowType"
      :min-date="minAvailableDate"
      @update:interval="updateInterval"
    />
  </v-layout>
</template>

<script>
import { QUICK_RANGES } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';

import { localQueryMixin } from '@/mixins/query/query';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';

import AvailabilityBar from '@/components/other/availability/partials/availability-bar.vue';
import ChartLoader from '@/components/widgets/chart/partials/chart-loader.vue';

export default {
  components: { ChartLoader, AvailabilityBar },
  mixins: [localQueryMixin, queryIntervalFilterMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    defaultTimeRange: {
      type: String,
      required: false,
    },
    defaultShowType: {
      type: Number,
      required: false,
    },
  },
  data() {
    const { start, stop } = QUICK_RANGES[this.defaultTimeRange];

    return {
      pending: false,
      availability: null,
      minDate: null,
      query: {
        interval: {
          from: start,
          to: stop,
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
    minDate() {
      const { from } = this.getIntervalQuery();

      if (this.minDate && from < this.minDate) {
        this.updateQueryField('interval', { ...this.query.interval, from: this.minDate });
      }
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    getQuery() {
      return {
        ...this.getIntervalQuery(),
      };
    },

    async fetchList() {
      this.pending = true;

      try {
        /**
         * TODO: Should be replaced on real fetch function
         */
        await new Promise(r => setTimeout(r, 2000));

        const minDate = new Date();
        minDate.setDate(minDate.getDate() - 3);

        this.minDate = Math.round(minDate.getTime() / 1000);

        this.availability = {
          uptime: Math.round(Math.random() * 100000),
          downtime: Math.round(Math.random() * 100000),
          inactive_time: Math.round(Math.random() * 1000),
        };
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
