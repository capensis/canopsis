<template lang="pug">
  v-layout
    v-flex(xs6)
      v-select(
        v-field="threshold.threshold_type",
        v-validate="'required'",
        :items="thresholdTypes",
        :label="$t('metaAlarmRule.thresholdType')",
        :error-messages="errors.collect('thresholdType')",
        name="thresholdType",
        item-text="label",
        item-value="field"
      )
    v-flex(xs6)
      v-text-field(
        v-if="threshold.threshold_type === $constants.META_ALARMS_THRESHOLD_TYPES.thresholdCount",
        v-field.number="threshold.threshold_count",
        v-validate="'required|numeric|min_value:0'",
        :label="$t('metaAlarmRule.thresholdCount')",
        :error-messages="errors.collect('thresholdCount')",
        :min="0",
        name="thresholdCount",
        type="number"
      )
      c-percents-field(
        v-else,
        v-field.number="threshold.threshold_rate",
        :label="$t('metaAlarmRule.thresholdRate')",
        name="thresholdRate"
      )
</template>

<script>
import { META_ALARMS_THRESHOLD_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'threshold',
    event: 'input',
  },
  props: {
    threshold: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    thresholdTypes() {
      return Object.values(META_ALARMS_THRESHOLD_TYPES).map(field => ({
        label: this.$t(`metaAlarmRule.${field}`),
        field,
      }));
    },
  },
};
</script>
