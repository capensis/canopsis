<template>
  <v-layout>
    <v-btn-toggle
      :value="value"
      :mandatory="mandatory"
      dense
      autofocus
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
import { computed } from 'vue';

import { ALARM_STATES } from '@/constants';

import { getAlarmStateColor } from '@/helpers/entities/alarm/color';

import { useI18n } from '@/hooks/i18n';

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
  setup(props) {
    const { t } = useI18n();

    const availableStates = computed(() => Object.entries(props.stateValues).map(([key, state]) => ({
      text: t(`modals.createChangeStateEvent.states.${key}`),
      state,
      color: getAlarmStateColor(state),
    })));

    return {
      availableStates,
    };
  },
};
</script>
