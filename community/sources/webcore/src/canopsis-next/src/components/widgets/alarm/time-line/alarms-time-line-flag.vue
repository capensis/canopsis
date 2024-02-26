<template>
  <v-badge
    :value="isActiveBadge"
    class="time-line-flag"
    color="transparent"
    overlap
  >
    <template #badge="">
      <v-icon
        class="time-line-flag__badge-icon"
        color="error"
        size="14"
      >
        error
      </v-icon>
    </template>
    <v-icon :style="{ color: style.color, caretColor: style.color }">
      {{ style.icon }}
    </v-icon>
  </v-badge>
</template>

<script>
import { ALARM_LIST_TIMELINE_STEPS } from '@/constants';

import { formatStep } from '@/helpers/entities/entity/formatting';

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
    error: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    style() {
      return formatStep(this.step);
    },

    isActiveBadge() {
      return [
        ALARM_LIST_TIMELINE_STEPS.declareTicketFail,
        ALARM_LIST_TIMELINE_STEPS.webhookFail,
      ].includes(this.step._t);
    },
  },
};
</script>

<style lang="scss">
.time-line-flag {
  &__badge-icon {
    background: white;
    border-radius: 50%;
  }
}
</style>
