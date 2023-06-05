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
    v-flex.ml-3(xs6)
      c-number-field(
        v-if="isThresholdCountType",
        v-field="threshold.threshold_count",
        :label="$t('metaAlarmRule.thresholdCount')",
        :min="0",
        name="thresholdCount"
      )
      c-percents-field(
        v-else,
        v-field="threshold.threshold_rate",
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
    isThresholdCountType() {
      return this.threshold.threshold_type === META_ALARMS_THRESHOLD_TYPES.thresholdCount;
    },

    thresholdTypes() {
      return Object.values(META_ALARMS_THRESHOLD_TYPES).map(field => ({
        label: this.$t(`metaAlarmRule.${field}`),
        field,
      }));
    },
  },
};
</script>
