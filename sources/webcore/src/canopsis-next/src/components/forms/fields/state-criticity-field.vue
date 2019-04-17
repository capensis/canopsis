<template lang="pug">
  v-layout(row)
    v-btn-toggle(
    :value="value",
    :mandatory="mandatory",
    @change="$emit('input', $event)"
    )
      v-btn(
      v-for="(state, key) in $constants.ENTITIES_STATES",
      :key="key",
      :value="state",
      :style="{ backgroundColor: $constants.ENTITIES_STATES_STYLES[state].color }",
      depressed
      ) {{ $t(`modals.createChangeStateEvent.states.${key}`) }}
</template>

<script>

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
  },
};
</script>
