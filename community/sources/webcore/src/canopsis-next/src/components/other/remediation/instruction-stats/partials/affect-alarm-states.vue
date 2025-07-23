<template>
  <div class="d-flex">
    <c-state-count-changes-chips
      :affect="alarmStates.minor"
      :state="ALARM_STATES.minor"
    />
    <c-state-count-changes-chips
      :affect="alarmStates.major"
      :state="ALARM_STATES.major"
      class="ml-2"
    />
    <c-state-count-changes-chips
      :affect="alarmStates.critical"
      :state="ALARM_STATES.critical"
      class="ml-2"
    />
  </div>
</template>

<script>
import { ALARM_STATES } from '@/constants';

import { getAlarmStateColor } from '@/helpers/entities/alarm/color';

export default {
  props: {
    alarmStates: {
      type: Object,
      required: true,
    },
  },
  setup() {
    return {
      ALARM_STATES,
    };
  },
  computed: {
    minorColor() {
      return getAlarmStateColor(ALARM_STATES.minor);
    },

    majorColor() {
      return getAlarmStateColor(ALARM_STATES.major);
    },

    criticalColor() {
      return getAlarmStateColor(ALARM_STATES.critical);
    },
  },
};
</script>
