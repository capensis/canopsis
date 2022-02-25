<template lang="pug">
  v-select(
    v-field="value",
    v-validate="rules",
    :items="actionTypes",
    :error-messages="errors.collect(name)",
    :label="label || $tc('common.condition')",
    :disabled="disabled",
    :name="name"
  )
</template>

<script>
import { FILTER_OPERATORS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: FILTER_OPERATORS.equal,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'condition',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    operators: {
      type: Array,
      default: () => Object.values(FILTER_OPERATORS),
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    actionTypes() {
      return this.operators.map(condition => ({
        value: condition,
        text: this.$t(`common.operators.${condition}`),
      }));
    },
  },
};
</script>
