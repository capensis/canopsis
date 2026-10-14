<template>
  <c-radio-group
    :model-value="modelValue"
    :name="name"
    :label="label"
    hide-details
    mandatory
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      v-for="type in types"
      :key="type.value"
      :label="type.label"
      :value="type.value"
      color="primary"
    />
  </c-radio-group>
</template>

<script>
import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

export default {
  props: {
    modelValue: {
      type: Number,
      default: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
    },
    name: {
      type: String,
      default: 'parameter',
    },
    label: {
      type: String,
      default: '',
    },
  },

  emits: ['update:modelValue'],

  computed: {
    types() {
      return [
        {
          value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
          label: this.$t('common.uptime'),
        },
        {
          value: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
          label: this.$t('common.downtime'),
        },
      ];
    },
  },
};
</script>
