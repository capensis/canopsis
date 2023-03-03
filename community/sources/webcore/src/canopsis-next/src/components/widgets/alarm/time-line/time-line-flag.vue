<template lang="pug">
  v-icon(:color="style.color") {{ style.icon }}
</template>

<script>
import { formatStep } from '@/helpers/formatting';

/**
 * Component for the flag on the alarms list's timeline
 *
 * @module alarm
 *
 * @prop {Object} step - step object
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
      return formatStep(this.step);
    },
  },
};
</script>
