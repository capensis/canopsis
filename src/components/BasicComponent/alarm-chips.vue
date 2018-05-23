<template lang="pug">
  v-chip(:color="color" text-color="white") {{ this.text }}
</template>

<script>
import { STATES_CHIPS_AND_FLAGS_STYLE, STATUS_CHIPS_AND_FLAGS_STYLE } from '@/config';

export default {
  name: 'state-chips',
  props: {
    val: {
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
      if (this.isStatus) {
        return STATUS_CHIPS_AND_FLAGS_STYLE[this.val].color;
      }
      return STATES_CHIPS_AND_FLAGS_STYLE[this.val].color;
    },
    text() {
      if (this.isStatus) {
        return STATUS_CHIPS_AND_FLAGS_STYLE[this.val].text;
      }
      return STATES_CHIPS_AND_FLAGS_STYLE[this.val].text;
    },
  },
};
</script>

<style scoped>

</style>
