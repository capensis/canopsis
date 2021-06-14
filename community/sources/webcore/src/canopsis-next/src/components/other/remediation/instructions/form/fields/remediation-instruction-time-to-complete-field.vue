<template lang="pug">
  v-layout
    v-flex(xs8)
      v-text-field(
        v-field.number="duration.value",
        v-validate="'required|min_value:1'",
        :label="$t('remediationInstructions.timeToComplete')",
        :error-messages="errors.collect(durationFieldName)",
        :min="0",
        :name="durationFieldName",
        type="number",
        box
      )
    v-flex.pl-3(xs4)
      v-select.time-complete-unit(
        v-field="duration.unit",
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
    prop: 'duration',
    event: 'input',
  },
  props: {
    duration: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: '',
    },
  },
  computed: {
    durationFieldName() {
      return `${this.name}.duration`;
    },

    unitFieldName() {
      return `${this.name}.unit`;
    },

    availableUnits() {
      return Object.values(omit(AVAILABLE_TIME_UNITS, ['year'])).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.duration.value),
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
