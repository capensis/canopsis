<template>
  <div class="d-flex">
    <c-state-count-changes-chips
      :affect="alarmStates.minor"
      :color="minorColor"
    />
    <c-state-count-changes-chips
      :affect="alarmStates.major"
      :color="majorColor"
      class="ml-2"
    />
    <c-state-count-changes-chips
      :affect="alarmStates.critical"
      :color="criticalColor"
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
