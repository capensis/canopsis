<template lang="pug">
  v-icon(large :color="color") {{ icon }}
</template>

<script>
import formatStateAndStatus from '@/helpers/state-and-status-formatting';

export default {
  name: 'alarm-flag',
  props: {
    value: {
      type: [Number, String],
      default: 0,
    },
    isStatus: {
      type: Boolean,
      default: false,
    },
    isCroppedState: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    color() {
      return formatStateAndStatus(this.value, this.isStatus, this.isCroppedState).color;
    },
    icon() {
      return formatStateAndStatus(this.value, this.isStatus, this.isCroppedState).icon;
    },
  },
};
</script>
