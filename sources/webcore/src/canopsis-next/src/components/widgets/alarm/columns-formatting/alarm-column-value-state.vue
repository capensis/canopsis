<template lang="pug">
  v-layout(align-center)
    alarm-chips(:type="$constants.ENTITY_INFOS_TYPE.state", :value="stateId")
    v-icon(v-if="showIcon", color="purple") account_circle
</template>

<script>
import { get } from 'lodash';

import AlarmChips from '../alarm-chips.vue';

/**
 * Component for the 'state' column of the alarms list
 *
 * @prop {Object} alarm - Object representing the alarm
 * @prop {String} propertyKey - Property name
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
  computed: {
    stateId() {
      return get(this.alarm, this.propertyKey);
    },
    showIcon() {
      return get(this.alarm, 'v.state._t') === this.$constants.EVENT_ENTITY_TYPES.changeState;
    },
  },
};
</script>
