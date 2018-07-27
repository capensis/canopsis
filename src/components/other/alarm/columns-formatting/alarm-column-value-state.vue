<template lang="pug">
  div
    alarm-chips(:type="ENTITY_INFOS_TYPE.state", :value="stateId")
    v-icon(:color="purple", v-if="showIcon") account_circle
</template>

<script>
import getProp from 'lodash/get';

import { ENTITY_INFOS_TYPE, EVENT_ENTITY_TYPES } from '@/constants';
import AlarmChips from '../alarm-chips.vue';

/**
 * Component for the 'state' column of the alarms list
 *
 * @module alarm
 *
 * @prop {Object} [alarm] - Object representing the alarm
 * @prop {String} [propertyKey] - Property name
 */
export default {
  components: { AlarmChips },
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
      ENTITY_INFOS_TYPE,
    };
  },
  computed: {
    stateId() {
      return getProp(this.alarm, this.propertyKey);
    },
    showIcon() {
      return getProp(this.alarm, 'v.state._t') === EVENT_ENTITY_TYPES.changeState;
    },
  },
};
</script>
