<template>
  <c-text-field
    :class="contentClass"
    :model-value="textValue"
    :label="label"
    :hide-details="hideDetails"
    disabled
    readonly
  />
</template>

<script>
import { convertDateToString } from '@/helpers/date/date';

export default {
  props: {
    value: {
      type: [String, Number, Date],
      default: null,
    },
    label: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
    },
    format: {
      type: String,
      default: 'short',
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    contentClass: {
      type: [String, Object],
      required: false,
    },
  },
  computed: {
    textValue() {
      return convertDateToString(this.value, this.format);
    },
  },
};
</script>
