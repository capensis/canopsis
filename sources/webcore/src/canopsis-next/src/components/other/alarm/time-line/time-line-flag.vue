<template lang="pug">
  div
    i.material-icons(:style="{color: style.color}") {{ style.icon }}
</template>

<script>
import { formatState, formatStatus, formatEvent } from '@/helpers/state-and-status-formatting';

/**
 * Component for the flag on the alarms list's timeline
 *
 * @module alarm
 *
 * @prop {Object} step - step object
 * @prop {Boolean} [cropped] - Boolean to determine if there's a cropped state or not
 */
export default {
  props: {
    step: {
      type: Object,
      required: true,
    },
  },
  computed: {
    style() {
      if (this.step._t.startsWith('status')) {
        return formatStatus(this.step.val);
      }

      if (this.step._t.startsWith('state')) {
        return formatState(this.step.val);
      }

      return formatEvent(this.step._t);
    },
  },
};
</script>
