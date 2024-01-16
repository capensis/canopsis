<template>
  <c-movable-card-iterator-field
    v-field="metrics"
    :addable="!onlyExternal"
    @add="addMetric"
  >
    <template #item="{ item, index }">
      <c-alarm-metric-preset-field
        v-field="metrics[index]"
        :with-color="withColor"
        :with-aggregate-function="withAggregateFunction"
        :parameters="parameters"
        :disabled-parameters="disabledParameters"
        :with-external="withExternal"
        :only-external="onlyExternal"
        :name="`${name}[${item.key}]`"
      />
    </template>
    <template #append="">
      <c-alert
        v-if="errorMessage"
        type="error"
      >
        {{ errorMessage }}
      </c-alert>
    </template>
    <template #actions="">
      <v-btn
        class="mr-2 mx-0"
        v-if="withExternal"
        color="primary"
        @click.prevent="addExternal"
      >
        {{ $t('kpi.addExternal') }}
      </v-btn>
      <v-btn
        class="mr-2 mx-0"
        v-if="withExternal"
        color="primary"
        @click.prevent="addAuto"
      >
        {{ $t('kpi.autoAdd') }}
      </v-btn>
    </template>
  </c-movable-card-iterator-field>
</template>

<script>
import { omit } from 'lodash';

import { ALARM_METRIC_PARAMETERS, AGGREGATE_FUNCTIONS } from '@/constants';

import { metricPresetToForm, isRatioMetric, isTimeMetric } from '@/helpers/entities/metric/form';

import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'metrics',
    event: 'input',
  },
  props: {
    metrics: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      default: 'metrics',
    },
    withColor: {
      type: Boolean,
      default: false,
    },
    withAggregateFunction: {
      type: Boolean,
      default: false,
    },
    withExternal: {
      type: Boolean,
      default: false,
    },
    parameters: {
      type: Array,
      default: () => Object.values(omit(ALARM_METRIC_PARAMETERS, ['timeToAck', 'timeToResolve'])),
    },
    onlyGroup: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: 1,
    },
    onlyExternal: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    errorMessage() {
      return this.errors.collect(this.name, undefined, false)
        ?.map(({ rule, msg }) => {
          const customMessage = {
            min_value: this.$t('kpi.errors.metricsMinLength', { count: this.min }),
          }[rule];

          return customMessage || msg;
        })
        .join('\n');
    },

    excludedParameters() {
      const [firstMetric] = this.metrics;

      if (!this.onlyGroup || !firstMetric?.metric) {
        return [];
      }

      if (isRatioMetric(firstMetric.metric)) {
        return this.parameters.filter(metric => !isRatioMetric(metric));
      }

      if (isTimeMetric(firstMetric.metric)) {
        return this.parameters.filter(metric => !isTimeMetric(metric));
      }

      return this.parameters.filter(metric => isTimeMetric(metric) || isRatioMetric(metric));
    },

    disabledParameters() {
      return [
        ...this.metrics.map(({ metric }) => metric),
        ...this.excludedParameters,
      ];
    },
  },
  watch: {
    metrics() {
      if (this.errorMessage) {
        this.$validator.validate(this.name);
      }
    },
  },
  created() {
    this.attachMinValueRule();
  },
  beforeDestroy() {
    this.detachRules();
  },
  methods: {
    addMetric(metric = {}) {
      this.addItemIntoArray(metricPresetToForm({
        aggregate_func: this.withAggregateFunction ? AGGREGATE_FUNCTIONS.avg : '',
        ...metric,
      }));
    },

    addExternal() {
      this.addMetric({ external: true, aggregate_func: AGGREGATE_FUNCTIONS.avg });
    },

    addAuto() {
      this.addMetric({ auto: true, aggregate_func: AGGREGATE_FUNCTIONS.avg });
    },

    attachMinValueRule() {
      this.$validator.attach({
        name: this.name,
        rules: { min_value: this.min },
        getter: () => (this.metrics.some(({ auto }) => auto) ? this.min : this.metrics.length),
        vm: this,
      });
    },

    detachRules() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
