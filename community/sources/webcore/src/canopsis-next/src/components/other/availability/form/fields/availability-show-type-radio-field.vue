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
import { AVAILABILITY_SHOW_TYPE } from '@/constants';

export default {
  props: {
    modelValue: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    name: {
      type: String,
      default: 'availability',
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
          value: AVAILABILITY_SHOW_TYPE.percent,
          label: this.$t('common.percent'),
        },
        {
          value: AVAILABILITY_SHOW_TYPE.duration,
          label: this.$t('common.duration'),
        },
      ];
    },
  },
};
</script>
