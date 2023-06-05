<template lang="pug">
  v-layout
    v-flex(xs8)
      c-number-field(
        v-field="duration.value",
        :label="$t('remediation.instruction.timeToComplete')",
        :name="durationFieldName",
        :disabled="disabled",
        :min="min",
        required,
        box
      )
    v-flex.pl-3(xs4)
      v-select.time-complete-unit(
        v-field="duration.unit",
        v-validate="'required'",
        :label="$t('common.unit')",
        :items="availableUnits",
        :name="unitFieldName",
        :disabled="disabled",
        hide-details
      )
</template>

<script>
import { SHORT_AVAILABLE_TIME_UNITS } from '@/constants';

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
    disabled: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: 1,
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
      return Object.values(SHORT_AVAILABLE_TIME_UNITS).map(({ value, text }) => ({
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
