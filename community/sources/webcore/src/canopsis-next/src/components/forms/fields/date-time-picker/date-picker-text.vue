<template>
  <v-text-field
    :class="contentClass"
    :value="value | date(format)"
    :label="label"
    :hide-details="hideDetails"
    disabled
    readonly
  />
</template>

<script>

export default {
  props: {
    value: {
      type: [String, Number],
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
};
</script>
