<template lang="pug">
  v-chip(:color="color" text-color="white") {{ this.text }}
</template>

<script>
import formatStateAndStatus from '@/services/state-and-status-formatting';

export default {
  name: 'alarm-chips',
  props: {
    value: {
      type: Number,
      default: 0,
    },
    isStatus: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    color() {
      return this.formatStateAndStatus(this.value, this.isStatus).color;
    },
    text() {
      return this.formatStateAndStatus(this.value, this.isStatus).text;
    },
  },
  methods: {
    formatStateAndStatus,
  },
};
</script>
