<template lang="pug">
  v-layout
    v-flex(xs6)
      v-text-field(
        v-field.number="delay.value",
        v-validate="fieldValidateRules",
        :label="$t('modals.createAction.fields.delay')",
        :error-messages="errors.collect('delayValue')",
        :min="0",
        name="delayValue",
        type="number"
      )
    v-flex(xs6)
      v-select(
        v-field="delay.unit",
        v-validate="unitValidateRules",
        :items="availableUnits",
        :label="$t('modals.createAction.fields.delayUnit')",
        :error-messages="errors.collect('delayUnit')",
        name="delayUnit",
        clearable
      )
</template>

<script>
import { isNumber } from 'lodash';
import { PERIODIC_REFRESH_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'delay',
    event: 'input',
  },
  props: {
    delay: {
      type: Object,
      default: () => ({}),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isRequired() {
      return this.required || isNumber(this.delay.value) || Boolean(this.delay.unit);
    },
    availableUnits() {
      return Object.values(PERIODIC_REFRESH_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, 2),
      }));
    },
    fieldValidateRules() {
      return {
        required: this.isRequired,
        numeric: true,
        min_value: 1,
      };
    },
    unitValidateRules() {
      return { required: this.isRequired };
    },
  },
};
</script>
