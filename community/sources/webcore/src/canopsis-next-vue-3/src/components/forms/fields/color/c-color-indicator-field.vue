<template>
  <v-radio-group
    :model-value="modelValue"
    :name="name"
    :label="label"
    :disabled="disabled"
    color="primary"
    v-bind="$attrs"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      :value="COLOR_INDICATOR_TYPES.state"
      :label="$t('settings.colorIndicator.fields.displayAsSeverity')"
      color="primary"
    />
    <v-radio
      :value="COLOR_INDICATOR_TYPES.impactState"
      :label="$t('settings.colorIndicator.fields.displayAsPriority')"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
      type: String,
      default: COLOR_INDICATOR_TYPES.state,
    },
    name: {
      type: String,
      default: 'colorIndicator',
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  data() {
    return {
      COLOR_INDICATOR_TYPES,
    };
  },
};
</script>
