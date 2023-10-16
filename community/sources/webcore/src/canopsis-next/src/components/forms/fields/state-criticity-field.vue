<template>
  <v-layout>
    <v-btn-toggle
      :value="value"
      :mandatory="mandatory"
      @change="$emit('input', $event)"
    >
      <v-btn
        v-for="{ color, value, text } in availableStates"
        :key="value"
        :value="value"
        :style="{ backgroundColor: color }"
        depressed="depressed"
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
        value: state,
        color: getEntityStateColor(state),
      }));
    },
  },
};
</script>
