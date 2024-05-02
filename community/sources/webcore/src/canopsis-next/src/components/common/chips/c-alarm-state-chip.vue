<template>
  <v-layout
    :class="classes"
    align-center
  >
    <c-alarm-chip
      :value="value"
      :badge-value="badgeValue"
      :small="small"
      @click="$emit('click', $event)"
    >
      <template v-if="appendIconName" #append>
        <v-icon color="white" size="14">
          {{ appendIconName }}
        </v-icon>
      </template>
    </c-alarm-chip>
    <v-icon
      v-if="showIcon"
      :size="small ? 14 : undefined"
      class="d-block"
      color="purple"
    >
      account_circle
    </v-icon>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { ALARM_LIST_STEPS, ALARM_STATES } from '@/constants';

/**
 * Component for the 'state' column of the alarms list
 *
 * @prop {Object} alarm - Object representing the alarm
 * @prop {String} propertyKey - Property name
 */
export default {
  props: {
    value: {
      type: [Number, String],
      default: ALARM_STATES.ok,
    },
    type: {
      type: String,
      default: '',
    },
    badgeValue: {
      type: [Number, String],
      required: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    appendIconName: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const showIcon = computed(() => props.type === ALARM_LIST_STEPS.changeState);
    const classes = computed(() => ({
      'c-alarm-state-chip--badge': props.badgeValue,
      'c-alarm-state-chip--small': props.small,
    }));

    return {
      showIcon,
      classes,
    };
  },
};
</script>

<style lang="scss">
.c-alarm-state-chip {
  &--badge {
    margin-top: 12px;
  }

  &--small {
    margin-top: 8px;
  }
}
</style>
