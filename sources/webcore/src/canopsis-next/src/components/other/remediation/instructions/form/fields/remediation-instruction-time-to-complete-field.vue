<template lang="pug">
  v-layout
    v-flex(xs8)
      v-text-field(
        v-field.number="value.interval",
        v-validate="'required|min_value:1'",
        :label="$t('remediationInstructions.timeToComplete')",
        :error-messages="timeToCompleteErrors",
        :min="0",
        :name="intervalFieldName",
        type="number",
        box
      )
    v-flex.pl-3(xs4)
      v-select.time-complete-unit(
        v-field="value.unit",
        v-validate="'required'",
        :items="availableUnits",
        :name="unitFieldName",
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
    name: {
      type: String,
      default: '',
    },
  },
  computed: {
    intervalFieldName() {
      return `${this.name}.interval`;
    },

    unitFieldName() {
      return `${this.name}.unit`;
    },

    timeToCompleteErrors() {
      return this.errors.collect(this.intervalFieldName)
        .map(error => error.replace(this.intervalFieldName, this.$t('remediationInstructions.timeToComplete')));
    },

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
