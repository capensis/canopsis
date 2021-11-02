<template lang="pug">
  div
    c-duration-field(v-field="form.duration", :units="snoozeUnits", required)
    v-layout(row)
      v-textarea(
        v-field="form.output",
        v-validate="isNoteRequired ? 'required' : ''",
        :label="$t('common.output')",
        :error-messages="errors.collect('output')",
        name="output"
      )
</template>

<script>
import { SNOOZE_DURATION_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    isNoteRequired: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    snoozeUnits() {
      return Object.values(SNOOZE_DURATION_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.form.duration.value || 0),
      }));
    },
  },
};
</script>
