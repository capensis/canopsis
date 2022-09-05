<template lang="pug">
  v-layout(:class="{ 'alarm-column-value-state--badge': badgeValue }", align-center)
    c-alarm-chip(:value="stateId", :badge-value="badgeValue")
    v-icon(v-if="showIcon", color="purple") account_circle
</template>

<script>
import { get } from 'lodash';

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
  },
  computed: {
    badgeValue() {
      return get(this.alarm, 'v.events_count');
    },

    stateId() {
      return get(this.alarm, this.propertyKey);
    },

    showIcon() {
      return get(this.alarm, 'v.state._t') === this.$constants.EVENT_ENTITY_TYPES.changeState;
    },
  },
};
</script>

<style lang="scss">
.alarm-column-value-state {
  &--badge {
    margin-top: 12px;
  }
}
</style>
