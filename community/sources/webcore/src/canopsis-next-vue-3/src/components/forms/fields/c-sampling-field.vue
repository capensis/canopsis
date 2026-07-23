<template>
  <c-select
    :model-value="value"
    :items="availableSamplings"
    :disabled="disabled"
    :label="label || $t('common.sampling')"
    :name="name"
    item-title="text"
    item-value="value"
    hide-details
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import { SAMPLINGS } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'sampling',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    samplings: {
      type: Array,
      default: () => Object.values(SAMPLINGS),
    },
  },

  emits: ['update:value'],

  computed: {
    availableSamplings() {
      return this.samplings.map((value) => ({
        value,
        text: this.$t(`common.times.${value}`, 2),
      }));
    },
  },
};
</script>
