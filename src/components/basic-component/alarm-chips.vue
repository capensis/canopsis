<template lang="pug">
  v-chip(:color="color" text-color="white") {{ this.text }}
</template>

<script>
import { ENTITY_STATES_STYLES, ENTITY_STATUS_STYLES } from '@/constants';

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
        return ENTITY_STATUS_STYLES[this.value].color;
      }
      return ENTITY_STATES_STYLES[this.value].color;
    },
    text() {
      if (this.isStatus) {
        return ENTITY_STATUS_STYLES[this.value].text;
      }
      return ENTITY_STATES_STYLES[this.value].text;
    },
  },
};
</script>
