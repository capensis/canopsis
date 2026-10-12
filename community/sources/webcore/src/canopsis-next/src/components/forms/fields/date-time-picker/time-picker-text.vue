<template>
  <c-text-field
    :model-value="value"
    :label="label"
    class="time-picker__select"
    placeholder="--:--"
    hide-details
    disabled
  />
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

export default {
  emits: ['update:value'],
  mixins: [formBaseMixin],
  props: {
    value: {
      type: String,
      default: null,
    },
    label: {
      type: String,
      default: null,
    },
  },
};
</script>
