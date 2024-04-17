<template>
  <v-layout>
    <v-btn-toggle
      :value="value"
      :mandatory="mandatory"
      dense
      @change="$emit('input', $event)"
    >
      <v-btn
        v-for="{ color, state, text } in availableStates"
        :key="state"
        :value="state"
        :style="{ backgroundColor: color }"
        depressed
      >
        {{ text }}
      </v-btn>
    </v-btn-toggle>
  </v-layout>
</template>

<script>
import { ALARM_STATES } from '@/constants';

import { getAlarmStateColor } from '@/helpers/entities/alarm/color';

export default {
  props: {
    value: {
      type: Number,
      default: null,
    },
    mandatory: {
      type: Boolean,
      default: false,
    },
    stateValues: {
      type: Object,
      default: () => ALARM_STATES,
    },
  },
  computed: {
    availableStates() {
      return Object.entries(this.stateValues).map(([key, state]) => ({
        text: this.$t(`modals.createChangeStateEvent.states.${key}`),
        state,
        color: getAlarmStateColor(state),
      }));
    },
  },
};
</script>
