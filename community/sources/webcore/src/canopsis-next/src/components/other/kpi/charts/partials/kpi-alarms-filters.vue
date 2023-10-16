<template>
  <div class="kpi-alarms-toolbar">
    <v-layout
      class="ml-4 my-4"
      wrap
    >
      <c-quick-date-interval-field
        class="mr-4"
        v-field="query.interval"
        :min="minFromTimestamp"
        :quick-ranges="quickRanges"
      />
      <c-sampling-field
        class="mr-4 kpi-alarms-toolbar__sampling"
        :value="query.sampling"
        @input="updateSampling"
      />
      <c-filter-field
        class="mr-4 kpi-alarms-toolbar__filters"
        v-field="query.filter"
      />
      <c-alarm-metric-parameters-field
        class="kpi-alarms-toolbar__parameters"
        v-field="query.parameters"
        hide-details
      />
    </v-layout>
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
import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { fromSeconds } from '@/helpers/date/duration';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'query',
    event: 'input',
  },
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

    updateSampling(sampling) {
      const { interval } = this.query;

      this.updateModel({
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
