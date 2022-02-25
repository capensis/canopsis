<template lang="pug">
  v-layout(row)
    v-btn-toggle(
      :value="value",
      :mandatory="mandatory",
      @change="$emit('input', $event)"
    )
      v-btn(
        :data-test="`stateCriticity-${state}`",
        v-for="(state, key) in stateValues",
        :key="key",
        :value="state",
        :style="{ backgroundColor: $constants.ENTITIES_STATES_STYLES[state].color }",
        depressed
      ) {{ $t(`modals.createChangeStateEvent.states.${key}`) }}
</template>

<script>
import { ENTITIES_STATES } from '@/constants';

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
};
</script>
