<template lang="pug">
  v-layout(data-test="durationField", row)
    v-flex(xs8)
      v-text-field(
        data-test="durationValue",
        v-field="value.duration",
        v-validate="'required|numeric|min_value:1'",
        :label="$t('modals.createSnoozeEvent.fields.duration')",
        :error-messages="errors.collect('duration')",
        data-vv-name="duration",
        type="number"
      )
    v-flex(data-test="durationType", xs4)
      v-select(
        v-field="value.durationType",
        v-validate="'required'",
        :items="availableUnits",
        :error-messages="errors.collect('durationType')",
        data-vv-name="durationType"
      )
</template>

<script>
import { DURATION_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    availableUnits() {
      return Object.values(DURATION_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, 2),
      }));
    },
  },
};
</script>
