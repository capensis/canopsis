<template>
  <c-select
    :model-value="value"
    :items="availableOperators"
    :label="label || $t('common.condition')"
    :disabled="disabled"
    :name="name"
    :required="required"
    item-title="text"
    item-value="value"
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import { PATTERN_OPERATORS } from '@/constants';

export default {
  emits: ['update:value'],
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
    availableOperators() {
      return this.operators.map(condition => ({
        value: condition,
        text: this.$t(`common.operators.${condition}`),
      }));
    },
  },
};
</script>
