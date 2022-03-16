<template lang="pug">
  v-select.c-pattern-operator-field(
    v-field="value",
    v-validate="rules",
    :items="availableOperators",
    :error-messages="errors.collect(name)",
    :label="label || $tc('common.condition')",
    :disabled="disabled",
    :name="name"
  )
    template(#selection="{ item }")
      span.ellipsis {{ item.text }}
</template>

<script>
import { PATTERN_OPERATORS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: PATTERN_OPERATORS.equal,
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
      default: () => Object.values(PATTERN_OPERATORS),
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    availableOperators() {
      return this.operators.map(condition => ({
        value: condition,
        text: this.$t(`common.operators.${condition}`),
      }));
    },
  },
};
</script>

<style lang="scss">
$selectIconWidth: 24px;

.c-pattern-operator-field {
  .v-select__selections {
    width: calc(100% - #{$selectIconWidth});
    flex-wrap: nowrap;
  }
}
</style>
