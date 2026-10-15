<template>
  <div class="kpi-alarms-toolbar">
    <c-layout class="ml-4 my-4" gap="3" wrap>
      <c-quick-date-interval-field
        :interval="query.interval"
        :min="minFromTimestamp"
        :quick-ranges="quickRanges"
        @update:interval="updateQueryField('interval', $event)"
      />
      <c-sampling-field
        :value="query.sampling"
        class="kpi-alarms-toolbar__sampling"
        @update:value="updateSampling"
      />
      <c-filter-field
        :value="query.filter"
        class="kpi-alarms-toolbar__filters"
        @update:value="updateQueryField('filter', $event)"
      />
      <kpi-alarm-metric-parameters-field
        :value="query.parameters"
        class="kpi-alarms-toolbar__parameters"
        hide-details
        @update:value="updateQueryField('parameters', $event)"
      />
    </c-layout>
  </div>
</template>

<script>
import {
  KPI_METRICS_MAX_ALARM_YEAR_INTERVAL_DIFF_IN_YEARS,
  METRICS_QUICK_RANGES,
  SAMPLINGS,
  TIME_UNITS,
} from '@/constants';

import { getDiffBetweenDates, subtractUnitFromDate } from '@/helpers/date/date';
import {
  convertStartDateIntervalToTimestamp,
  convertStopDateIntervalToTimestamp,
} from '@/helpers/date/date-intervals';
import { fromSeconds } from '@/helpers/date/duration';

import KpiAlarmMetricParametersField from '../form/fields/kpi-alarm-metric-parameters-field.vue';

export default {
  components: { KpiAlarmMetricParametersField },

  props: {
    query: {
      type: Object,
      required: true,
    },
    minDate: {
      type: Number,
      required: false,
    },
  },

  emits: ['update:query'],

  computed: {
    quickRanges() {
      return Object.values(METRICS_QUICK_RANGES);
    },

    minFromTimestamp() {
      if (this.isHourSampling(this.query.sampling)) {
        const yearAgoDate = this.getDateWithMaxIntervalDiff(
          convertStopDateIntervalToTimestamp(this.query.interval.to),
        );

        return this.minDate > yearAgoDate ? this.minDate : yearAgoDate;
      }

      return this.minDate;
    },
  },

  methods: {
    isHourSampling(sampling) {
      return sampling === SAMPLINGS.hour;
    },

    getDateWithMaxIntervalDiff(to) {
      return subtractUnitFromDate(
        to,
        KPI_METRICS_MAX_ALARM_YEAR_INTERVAL_DIFF_IN_YEARS,
        TIME_UNITS.year,
      );
    },

    getIntervalWithMaxYearRange(interval) {
      const { from, to } = interval;

      const toTimestamp = convertStopDateIntervalToTimestamp(to);
      const fromTimestamp = convertStartDateIntervalToTimestamp(from);

      const secondsDiff = getDiffBetweenDates(toTimestamp, fromTimestamp);
      const diffInYears = fromSeconds(secondsDiff, TIME_UNITS.year);

      return diffInYears > 1
        ? { to, from: this.getDateWithMaxIntervalDiff(toTimestamp) }
        : interval;
    },

    updateQueryField(field, value) {
      this.$emit('update:query', {
        ...this.query,
        [field]: value,
      });
    },

    updateSampling(sampling) {
      const { interval } = this.query;

      this.$emit('update:query', {
        ...this.query,
        interval: this.isHourSampling(sampling)
          ? this.getIntervalWithMaxYearRange(interval)
          : interval,
        sampling,
      });
    },
  },
};
</script>

<style scoped lang="scss">
.kpi-alarms-toolbar {
  &__sampling {
    max-width: 200px;
  }

  &__filters {
    max-width: 200px;
  }

  &__parameters {
    max-width: 300px;
  }
}
</style>
