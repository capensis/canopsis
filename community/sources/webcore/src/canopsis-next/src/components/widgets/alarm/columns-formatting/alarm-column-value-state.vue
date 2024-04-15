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
import { get } from 'lodash';
import { computed } from 'vue';

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
    appendIconName: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const badgeValue = computed(() => get(props.alarm, 'v.events_count'));
    const stateId = computed(() => get(props.alarm, props.propertyKey));
    const showIcon = computed(() => get(props.alarm, 'v.state._t') === ALARM_LIST_ACTIONS_TYPES.changeState);
    const classes = computed(() => ({
      'alarm-column-value-state--badge': badgeValue.value,
      'alarm-column-value-state--small': props.small,
    }));

    return {
      badgeValue,
      stateId,
      showIcon,
      classes,
    };
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
