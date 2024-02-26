<template>
  <v-layout
    :class="classes"
    align-center
  >
    <c-alarm-chip
      :value="stateId"
      :badge-value="badgeValue"
      :small="small"
      @click="$emit('click', $event)"
    />
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
import { get } from 'lodash';

import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

/**
 * Component for the 'state' column of the alarms list
 *
 * @prop {Object} alarm - Object representing the alarm
 * @prop {String} propertyKey - Property name
 */
export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    propertyKey: {
      type: String,
      default: 'v.state.val',
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    classes() {
      return {
        'alarm-column-value-state--badge': this.badgeValue,
        'alarm-column-value-state--small': this.small,
      };
    },

    badgeValue() {
      return get(this.alarm, 'v.events_count');
    },

    stateId() {
      return get(this.alarm, this.propertyKey);
    },

    showIcon() {
      return get(this.alarm, 'v.state._t') === ALARM_LIST_ACTIONS_TYPES.changeState;
    },
  },
};
</script>

<style lang="scss">
.alarm-column-value-state {
  &--badge {
    margin-top: 12px;
  }

  &--small {
    margin-top: 8px;
  }
}
</style>
