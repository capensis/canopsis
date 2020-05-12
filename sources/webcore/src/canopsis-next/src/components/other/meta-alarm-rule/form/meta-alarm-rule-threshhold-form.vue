<template lang="pug">
  v-layout
    v-flex(xs6)
      v-text-field(
        v-field.number="threshold.threshold_rate",
        v-validate="'required|decimal|min_value:0|max_value:1'",
        :label="$t('metaAlarmRule.fields.thresholdRate')",
        :error-messages="errors.collect('thresholdRate')",
        :min="0",
        :max="1",
        :step="0.1",
        name="thresholdRate",
        type="number"
      )
    v-flex(xs6)
      v-text-field(
        v-field.number="threshold.threshold_count",
        v-validate="'required|numeric|min_value:1'",
        :label="$t('metaAlarmRule.fields.thresholdCount')",
        :error-messages="errors.collect('thresholdCount')",
        :min="0",
        name="thresholdCount",
        type="number"
      )
</template>

<script>
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
};
</script>

