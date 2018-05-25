<template lang="pug">
  v-icon(large :color="color") {{ icon }}
</template>

<script>
import { STATES_CHIPS_AND_FLAGS_STYLE, STATUS_CHIPS_AND_FLAGS_STYLE } from '@/config';

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
      if (this.isStatus) {
        return STATUS_CHIPS_AND_FLAGS_STYLE[this.value].color;
      } else if (this.isCroppedState) {
        return 'black';
      }
      return STATES_CHIPS_AND_FLAGS_STYLE[this.value].color;
    },
    icon() {
      if (this.isStatus) {
        return STATUS_CHIPS_AND_FLAGS_STYLE[this.value].icon;
      } else if (this.isCroppedState) {
        return 'vertical_align_center';
      }
      return STATES_CHIPS_AND_FLAGS_STYLE[this.value].icon;
    },
  },
};
</script>
