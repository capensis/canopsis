<template lang="pug">
  v-icon(large :color="color") {{ icon }}
</template>

<script>
import { ENTITY_STATES_STYLES, ENTITY_STATUS_STYLES } from '@/constants';

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
        return ENTITY_STATUS_STYLES[this.value].color;
      } else if (this.isCroppedState) {
        return 'black';
      }
      return ENTITY_STATES_STYLES[this.value].color;
    },
    icon() {
      if (this.isStatus) {
        return ENTITY_STATUS_STYLES[this.value].icon;
      } else if (this.isCroppedState) {
        return 'vertical_align_center';
      }
      return ENTITY_STATES_STYLES[this.value].icon;
    },
  },
};
</script>
