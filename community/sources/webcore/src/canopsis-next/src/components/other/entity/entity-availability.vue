<template>
  <c-layout class="entity-availability position-relative" align="start" justify="center">
    <c-progress-overlay :pending="pending" />
    <availability-bar
      v-if="availability"
      :query="query"
      :availability="availability"
      :default-show-type="defaultShowType"
      class="entity-availability__content"
      @update:interval="updateInterval"
    />
  </c-layout>
</template>

<script>
import { mapActions } from 'pinia';

import { QUICK_RANGES } from '@/constants';

import { useMetricsStore } from '@/stores/modules/metrics';

import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { entitiesInfoMixin } from '@/mixins/entities/info';

import AvailabilityBar from '@/components/other/availability/partials/availability-bar.vue';

export default {
  inject: ['$periodicRefresh'],

  components: { AvailabilityBar },

  mixins: [entitiesInfoMixin],

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
    const { start, stop } = QUICK_RANGES[this.defaultTimeRange] ?? QUICK_RANGES.today;

    return {
      pending: false,
      availability: null,
      query: {
        interval: {
          from: start,
          to: stop,
        },
      },
    };
  },

  mounted() {
    this.fetchList();
    this.$periodicRefresh?.register?.(this.fetchList);
  },

  beforeUnmount() {
    this.$periodicRefresh?.unregister?.(this.fetchList);
  },

  methods: {
    ...mapActions(useMetricsStore, ['fetchAvailabilityWithoutStore']),

    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
      this.fetchList();
    },

    getIntervalQuery() {
      const { interval } = this.query;

      if (!interval) {
        return {};
      }

      return convertMetricIntervalToTimestamp({
        interval,
        timezone: this.timezone,
      });
    },

    getQuery() {
      return {
        _id: this.entity._id,
        ...this.getIntervalQuery(),
      };
    },

    async fetchList() {
      this.pending = true;

      try {
        this.availability = await this.fetchAvailabilityWithoutStore({
          params: this.getQuery(),
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.entity-availability {
  min-height: 100px;

  &__content {
    width: 100%;
    max-width: 900px;
  }
}
</style>
