<template>
  <c-select
    :model-value="currentValue"
    :items="actionTypes"
    :label="label || $t('common.type')"
    :disabled="disabled"
    :name="name"
    :required="required"
    :multiple="isMultiple"
    :chips="isMultiple"
    :closable-chips="isMultiple"
    item-title="text"
    item-value="value"
    @update:model-value="updateValue"
  />
</template>

<script>
import { isArray } from 'lodash';

import { BASIC_ENTITY_TYPES } from '@/constants';

export default {
  emits: ['update:value', 'update:modelValue'],

  props: {
    value: {
      type: [String, Array],
      required: false,
    },
    modelValue: {
      type: [String, Array],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    types: {
      type: Array,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    currentValue() {
      return this.modelValue ?? this.value;
    },

    actionTypes() {
      const types = this.types ?? Object.values(BASIC_ENTITY_TYPES);

      return types.map(type => ({
        value: type,
        text: this.$t(`entity.types.${type}`),
      }));
    },

    isMultiple() {
      return isArray(this.currentValue);
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
