<template>
  <v-layout
    :style="itemStyle"
    class="numbers-metrics-item"
    column
    align-center
  >
    <v-layout align-center>
      <div class="numbers-metrics-item__value">
        {{ value }}
      </div>
      <c-help-icon
        v-if="showTrend && trendEnabled"
        :text="trendTooltipText"
        :icon-class="{ 'numbers-metrics-item__trend': true, 'numbers-metrics-item__trend--up': trendUp }"
        icon="arrow_downward"
        size="0.75em"
        top
      />
    </v-layout>
    <div class="numbers-metrics-item__title">
      {{ title }}
    </div>
  </v-layout>
</template>

<script>
import { isUndefined, lowerFirst } from 'lodash';

import { AGGREGATE_FUNCTIONS, DATETIME_FORMATS, NUMBERS_CHART_DEFAULT_FONT_SIZE } from '@/constants';

import { convertMetricValueByUnit, convertMetricValueToString } from '@/helpers/entities/metric/list';
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
    valueFontSize: {
      type: Number,
      default: NUMBERS_CHART_DEFAULT_FONT_SIZE,
    },
  },
  computed: {
    itemStyle() {
      return {
        fontSize: `${this.valueFontSize}px`,
      };
    },

    value() {
      const { value, unit, title: metric } = this.metric;

      const preparedValue = convertMetricValueByUnit(value, unit);

      return convertMetricValueToString({ value: preparedValue, metric, unit });
    },

    title() {
      if (this.metric.label) {
        return this.metric.label;
      }

      const messageKey = `alarm.metrics.${this.metric.title}`;
      const title = this.$te(messageKey) ? this.$t(messageKey) : this.metric.title;

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

    trendEnabled() {
      return this.trendExist && this.metric.value !== this.metric.previous_metric;
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
    line-height: normal;
    text-align: center;
  }

  &__trend {
    transform: rotate(0deg);
    transition: transform linear .3s;

    &--up {
      transform: rotate(180deg);
    }
  }

  &__title {
    font-size: 0.2em;
  }
}
</style>
