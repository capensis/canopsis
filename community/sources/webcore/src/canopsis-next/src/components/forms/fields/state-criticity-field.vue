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
import { ENTITIES_STATES } from '@/constants';

import { getEntityStateColor } from '@/helpers/entities/entity/color';

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
      default: () => ENTITIES_STATES,
    },
  },
  computed: {
    availableStates() {
      return Object.entries(this.stateValues).map(([key, state]) => ({
        text: this.$t(`modals.createChangeStateEvent.states.${key}`),
        state,
        color: getEntityStateColor(state),
      }));
    },
  },
};
</script>
