<template lang="pug">
  v-chip(:color="color" text-color="white") {{ this.text }}
</template>

<script>
import { STATES_CHIPS_AND_FLAGS_STYLE, STATUS_CHIPS_AND_FLAGS_STYLE } from '@/config';

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
      if (this.isStatus) {
        return STATUS_CHIPS_AND_FLAGS_STYLE[this.value].color;
      }
      return STATES_CHIPS_AND_FLAGS_STYLE[this.value].color;
    },
    text() {
      if (this.isStatus) {
        return STATUS_CHIPS_AND_FLAGS_STYLE[this.value].text;
      }
      return STATES_CHIPS_AND_FLAGS_STYLE[this.value].text;
    },
  },
};
</script>
