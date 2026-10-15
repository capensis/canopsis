<template>
  <c-layout>
    <c-state-count-changes-chips :affect="alarmStates.minor" :state="ALARM_STATES.minor" />
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
  </c-layout>
</template>

<script>
import { ALARM_STATES } from '@/constants';

export default {
  props: {
    alarmStates: {
      type: Object,
      required: true,
    },
  },

  setup() {
    return { ALARM_STATES };
  },
};
</script>
