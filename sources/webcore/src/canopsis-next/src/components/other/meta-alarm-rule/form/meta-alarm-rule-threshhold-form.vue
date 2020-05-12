<template lang="pug">
  v-layout
    v-flex(xs6)
      percents-field(
        v-field.number="threshold.threshold_rate",
        :label="$t('metaAlarmRule.fields.thresholdRate')",
        name="thresholdRate"
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
import PercentsField from '@/components/forms/fields/percents.vue';

export default {
  components: { PercentsField },
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

