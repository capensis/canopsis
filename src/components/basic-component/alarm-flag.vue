<template lang="pug">
  v-icon(large :color="color") {{ icon }}
</template>

<script>
import formatStateAndStatus from '@/services/state-and-status-formatting';

export default {
  name: 'alarm-flag',
  props: {
    value: {
      type: Number,
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
      return this.formatStateAndStatus(this.value, this.isStatus, this.isCroppedState).color;
    },
    icon() {
      return this.formatStateAndStatus(this.value, this.isStatus, this.isCroppedState).icon;
    },
  },
  methods: {
    formatStateAndStatus,
  },
};
</script>
