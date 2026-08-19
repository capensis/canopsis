<template>
  <widget-settings-item :title="$t('settings.chart.selectMetrics')">
    <c-layout gap="3" column>
      <v-card v-for="(metric, index) in modelValue" :key="metric.key" class="lighten-2">
        <c-layout class="pt-2 px-2" align="center" justify="space-between">
          <c-layout>
            <v-btn
              :disabled="index === 0"
              icon="arrow_upward"
              variant="text"
              size="small"
              @click.prevent="moveMetric(index, -1)"
            />
            <v-btn
              :disabled="index === modelValue.length - 1"
              icon="arrow_downward"
              variant="text"
              size="small"
              @click.prevent="moveMetric(index, 1)"
            />
          </c-layout>
          <v-btn
            :disabled="modelValue.length <= min"
            icon="close"
            variant="text"
            size="small"
            @click.prevent="removeMetric(index)"
          />
        </c-layout>
        <v-card-text class="px-3 pb-3">
          <chart-metric-preset
            :model-value="metric"
            :name="`metric-${index}`"
            :with-color="withColor"
            :with-aggregate-function="withAggregateFunction"
            :parameters="parameters"
            :disabled-parameters="disabledParameters"
            @update:model-value="updatePreset(index, $event)"
          />
        </v-card-text>
      </v-card>

      <c-layout gap="2" justify="end">
        <v-btn
          v-if="withExternal"
          color="primary"
          variant="outlined"
          size="small"
          @click.prevent="addMetric({ external: true })"
        >
          {{ $t('kpi.addExternal') }}
        </v-btn>
        <v-btn
          v-if="withExternal"
          color="primary"
          variant="outlined"
          size="small"
          @click.prevent="addMetric({ auto: true, external: true })"
        >
          {{ $t('kpi.autoAdd') }}
        </v-btn>
        <v-btn
          v-if="!onlyExternal"
          color="primary"
          variant="outlined"
          size="small"
          @click.prevent="addMetric()"
        >
          {{ $t('common.add') }}
        </v-btn>
      </c-layout>
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { omit } from 'lodash';

import { AGGREGATE_FUNCTIONS, ALARM_METRIC_PARAMETERS } from '@/constants';

import { metricPresetToForm, isRatioMetric, isTimeMetric } from '@/helpers/entities/metric/form';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

import ChartMetricPreset from './chart-metric-preset.vue';

export default {
  components: { ChartMetricPreset, WidgetSettingsItem },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    withColor: {
      type: Boolean,
      default: false,
    },
    withExternal: {
      type: Boolean,
      default: false,
    },
    withAggregateFunction: {
      type: Boolean,
      default: false,
    },
    parameters: {
      type: Array,
      default: () => Object.values(omit(ALARM_METRIC_PARAMETERS, ['timeToAck', 'timeToResolve'])),
    },
    min: {
      type: Number,
      default: 0,
    },
    onlyGroup: {
      type: Boolean,
      default: false,
    },
    onlyExternal: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    disabledParameters() {
      const metrics = this.modelValue.map(({ metric }) => metric);
      const firstMetric = this.modelValue[0]?.metric;

      if (!this.onlyGroup || !firstMetric) {
        return metrics;
      }

      const excluded = isRatioMetric(firstMetric)
        ? this.parameters.filter(metric => !isRatioMetric(metric))
        : isTimeMetric(firstMetric)
          ? this.parameters.filter(metric => !isTimeMetric(metric))
          : this.parameters.filter(metric => isTimeMetric(metric) || isRatioMetric(metric));

      return [...metrics, ...excluded];
    },
  },

  methods: {
    addMetric(extra = {}) {
      this.$emit('update:modelValue', [
        ...this.modelValue,
        metricPresetToForm({
          aggregate_func:
            this.withAggregateFunction || extra.external || extra.auto
              ? AGGREGATE_FUNCTIONS.avg
              : '',
          ...extra,
        }),
      ]);
    },

    updatePreset(index, value) {
      const metrics = [...this.modelValue];
      metrics[index] = value;
      this.$emit('update:modelValue', metrics);
    },

    moveMetric(index, offset) {
      const targetIndex = index + offset;

      if (targetIndex < 0 || targetIndex >= this.modelValue.length) {
        return;
      }

      const metrics = [...this.modelValue];
      [metrics[index], metrics[targetIndex]] = [metrics[targetIndex], metrics[index]];
      this.$emit('update:modelValue', metrics);
    },

    removeMetric(index) {
      this.$emit(
        'update:modelValue',
        this.modelValue.filter((_, itemIndex) => itemIndex !== index),
      );
    },
  },
};
</script>
