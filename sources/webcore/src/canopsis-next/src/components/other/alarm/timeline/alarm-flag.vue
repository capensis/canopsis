<template lang="pug">
  v-icon(large :color="color") {{ icon }}
</template>

<script>
import formatStateAndStatus from '@/helpers/state-and-status-formatting';

/**
 * Component for the flag on the alarms list's timeline
 *
 * @module alarm
 *
 * @prop {Number,String} [value] - Value of the state/status
 * @prop {Boolean} [isStatus] - Boolean to determine if this is a state, or a status
 * @prop {Boolean} [isCroppedState] - Boolean to determine if there's a cropped state or not
 */
export default {
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
