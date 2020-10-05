<template lang="pug">
  v-layout
    v-flex(xs8)
      v-text-field(
        v-field.number="value.interval",
        v-validate="'required|min_value:1'",
        :label="$t('remediationInstructions.timeToComplete')",
        :error-messages="errors.collect('interval')",
        :disabled="disabled",
        :min="0",
        name="interval",
        type="number",
        box
      )
    v-flex.pl-3(xs4)
      v-select.time-complete-unit(
        v-field="value.unit",
        v-validate="'required'",
        :items="availableUnits",
        :disabled="disabled",
        name="unit",
        hide-details
      )
</template>

<script>
import { omit } from 'lodash';

import { AVAILABLE_TIME_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableUnits() {
      return Object.values(omit(AVAILABLE_TIME_UNITS, ['year'])).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.value.interval),
      }));
    },
  },
};
</script>

<style lang="scss">
  .time-complete-unit .v-input__slot {
    &:before, &:after {
      content: none !important;
    }
  }
</style>
