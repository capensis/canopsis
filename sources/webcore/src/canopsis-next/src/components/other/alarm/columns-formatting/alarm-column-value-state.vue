<template lang="pug">
  div
    span(:class="[`bg-${color}`, 'badge']") {{ text }}
    v-icon(:color="color", v-if="showIcon") account_circle
</template>

<script>
import getProp from 'lodash/get';

import { ENTITIES_STATES, EVENT_ENTITY_TYPES } from '@/constants';

/**
 * Component for the 'state' column of the alarms list
 *
 * @module alarm
 *
 * @prop {Object} [alarm] - Object representing the alarm
 * @prop {String} [propertyKey] - Property name
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
  data() {
    return {
      colors: {
        [ENTITIES_STATES.ok]: 'green',
        [ENTITIES_STATES.minor]: 'yellow',
        [ENTITIES_STATES.major]: 'orange',
        [ENTITIES_STATES.critical]: 'red',
      },
    };
  },
  computed: {
    stateId() {
      return getProp(this.alarm, this.propertyKey);
    },
    showIcon() {
      return getProp(this.alarm, 'v.state._t') === EVENT_ENTITY_TYPES.changeState;
    },
    color() {
      return this.colors[this.stateId] || 'purple';
    },
    text() {
      return this.$t(`tables.alarmStates.${this.stateId}`) || 'Unknown';
    },
  },
};
</script>

<style scoped>
  .badge {
    display: inline-block;
    min-width: 10px;
    padding: 3px 7px;
    font-size: 12px;
    font-weight: 700;
    line-height: 1;
    color: #fff;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    background-color: #777;
    border-radius: 10px;
  }

  .badge.bg-green {
    background-color: #00a65a
  }

  .badge.bg-red {
    background-color: #CF0000
  }

  .badge.bg-yellow {
    background-color: #FFE000
  }

  .badge.bg-orange {
    background-color: #FF9900
  }
</style>
