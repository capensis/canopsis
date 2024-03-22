<template>
  <v-layout class="entity-availability position-relative" align-start justify-center>
    <c-progress-overlay :pending="pending" />
    <availability-bar
      v-if="availability"
      :query="query"
      :availability="availability"
      :default-show-type="defaultShowType"
      :min-date="minAvailableDate"
      class="entity-availability__content"
      @update:interval="updateInterval"
    />
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { QUICK_RANGES } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';

import { localQueryMixin } from '@/mixins/query/query';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';

import AvailabilityBar from '@/components/other/availability/partials/availability-bar.vue';

const { mapActions: mapAvailabilityActions } = createNamespacedHelpers('availability');

export default {
  components: { AvailabilityBar },
  mixins: [localQueryMixin, queryIntervalFilterMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    defaultTimeRange: {
      type: String,
      default: QUICK_RANGES.today.value,
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
    ...mapAvailabilityActions({
      fetchAvailabilityWithoutStore: 'fetchAvailabilityWithoutStore',
    }),

    customQueryCondition(query, oldQuery) {
      return isMetricsQueryChanged(query, oldQuery, this.minDate);
    },

    getQuery() {
      return {
        ...this.getIntervalQuery(),
      };
    },

    async fetchList() {
      this.pending = true;

      try {
        const { min_date: minDate, availability } = await this.fetchAvailabilityWithoutStore({
          id: this.entity._id,
          params: this.getQuery(),
        });

        this.minDate = minDate;
        this.availability = availability;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss">
.entity-availability {
  min-height: 100px;

  &__content {
    max-width: 900px;
    width: 100%;
  }
}
</style>
