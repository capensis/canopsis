<template lang="pug">
  v-chip(:color="color" text-color="white") {{ this.text }}
</template>

<script>
import formatStateAndStatus from '@/helpers/state-and-status-formatting';

export default {
  name: 'alarm-chips',
  props: {
    value: {
      type: [Number, String],
      default: 0,
    },
    isStatus: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    color() {
      return formatStateAndStatus(this.value, this.isStatus).color;
    },
    text() {
      return formatStateAndStatus(this.value, this.isStatus).text;
    },
  },
};
</script>
