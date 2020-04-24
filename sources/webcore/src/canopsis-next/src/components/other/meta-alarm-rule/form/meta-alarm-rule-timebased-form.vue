<template lang="pug">
  v-layout
    v-flex(xs12)
      v-text-field(
        v-field.number="timebased.time_interval",
        v-validate="'required|numeric|min_value:1'",
        :label="$t('metaAlarmRule.fields.timeInterval')",
        :error-messages="errors.collect('timeInterval')",
        :min="1",
        name="timeInterval",
        type="number"
      )
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'timebased',
    event: 'input',
  },
  props: {
    timebased: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

