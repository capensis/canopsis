<template>
  <c-select
    :model-value="value"
    :label="label || $t('common.states')"
    :items="availableStateTypes"
    :disabled="disabled"
    :name="name"
    :required="required"
    item-title="text"
    item-value="value"
    hide-details
    @update:model-value="updateValue"
  />
</template>

<script>
import { ALARM_STATES } from '@/constants';

export default {
  emits: ['update:value', 'update:modelValue'],
  props: {
    value: {
      type: Number,
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'state',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableStateTypes() {
      return Object.values(ALARM_STATES).map(value => ({
        value,
        text: this.$t(`common.stateTypes.${value}`),
      }));
    },
  },
  methods: {
    updateValue(value) {
      this.$emit('update:value', value);
      this.$emit('update:modelValue', value);
    },
  },
};
</script>
