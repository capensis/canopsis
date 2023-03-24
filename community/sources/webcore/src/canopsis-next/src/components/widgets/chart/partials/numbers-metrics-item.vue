<template lang="pug">
  v-layout.numbers-metrics-item(column, align-center)
    v-layout(row, align-center)
      div.numbers-metrics-item__value {{ value }}
      c-help-icon(
        v-if="showTrend && trendExist",
        :text="trendTooltipText",
        :icon-class="{ 'numbers-metrics-item__trend--up': trendUp }",
        icon="arrow_downward",
        size="84",
        top
      )
    div.numbers-metrics-item__title {{ title }}
</template>

<script>
import { isUndefined, lowerFirst } from 'lodash';

import { AGGREGATE_FUNCTIONS, DATETIME_FORMATS } from '@/constants';

import { isTimeMetric } from '@/helpers/metrics';
import { convertDurationToString } from '@/helpers/date/duration';
import { convertDateToTimezoneDateString } from '@/helpers/date/date';

export default {
  inject: ['$system'],
  props: {
    metric: {
      type: Object,
      required: true,
    },
    showTrend: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    value() {
      if (isTimeMetric(this.metric.title)) {
        return convertDurationToString(this.metric.value);
      }

      return this.metric.value;
    },

    title() {
      const title = this.$t(`alarm.metrics.${this.metric.title}`);

      if (!this.metric.aggregate_func) {
        return title;
      }

      const prefix = this.metric.aggregate_func === AGGREGATE_FUNCTIONS.sum
        ? this.$t('common.total')
        : this.$t(`kpi.aggregateFunctions.${this.metric.aggregate_func}`);

      return `${prefix} ${lowerFirst(title)}`;
    },

    trendExist() {
      return !isUndefined(this.metric.previous_metric);
    },

    trendUp() {
      return this.metric.value > this.metric.previous_metric;
    },

    trendTooltipText() {
      const { previous_interval: previousInterval } = this.metric;

      return this.$t('kpi.periodTrend', {
        count: this.metric.previous_metric,
        from: convertDateToTimezoneDateString(previousInterval.from, this.$system.timezone, DATETIME_FORMATS.short),
        to: convertDateToTimezoneDateString(previousInterval.to, this.$system.timezone, DATETIME_FORMATS.short),
      });
    },
  },
};
</script>

<style lang="scss">
.numbers-metrics-item {
  &__value {
    font-size: 110px;
    text-align: center;
  }

  &__trend {
    transform: rotate(0deg);
    transition: linear .3s;

    &--up {
      transform: rotate(180deg);
    }
  }

  &__title {
    font-size: 18px;
  }
}
</style>
