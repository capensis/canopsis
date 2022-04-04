<template lang="pug">
  div
    v-tooltip(v-for="action in actions", :key="action.type", top)
      template(#activator="{ on }")
        service-entity-alarm-instruction-menu(
          v-if="action.type === $constants.WEATHER_ACTIONS_TYPES.executeInstruction",
          v-on="on",
          :assigned-instructions="assignedInstructions",
          :icon="action.icon",
          @execute="$listeners.execute"
        )
        v-btn(
          v-else,
          v-on="on",
          :disabled="action.disabled",
          depressed,
          small,
          light,
          @click.stop="$emit('apply', action)"
        )
          v-icon {{ action.icon }}
      span {{ $t(`serviceWeather.actions.${action.type}`) }}
</template>

<script>
import ServiceEntityAlarmInstructionMenu from './service-entity-alarm-instruction-menu.vue';

export default {
  components: {
    ServiceEntityAlarmInstructionMenu,
  },
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    assignedInstructions: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
