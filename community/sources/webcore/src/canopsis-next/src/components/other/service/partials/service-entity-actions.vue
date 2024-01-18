<template>
  <v-layout wrap>
    <v-tooltip
      v-for="action in actions"
      :key="action.type"
      top
    >
      <template #activator="{ on }">
        <span v-on="on">
          <service-entity-alarm-instruction-menu
            v-if="action.type === $constants.WEATHER_ACTIONS_TYPES.executeInstruction"
            :icon="action.icon"
            :entity="entity"
            :assigned-instructions="assignedInstructions"
            @execute="$listeners.execute"
          />
          <v-btn
            v-else
            class="ml-2"
            :disabled="action.disabled"
            :loading="action.loading"
            depressed
            small
            light
            @click.stop="$emit('apply', action)"
          >
            <v-icon>{{ action.icon }}</v-icon>
          </v-btn>
        </span>
      </template>
      <span>{{ $t(`serviceWeather.actions.${action.type}`) }}</span>
    </v-tooltip>
  </v-layout>
</template>

<script>
import ServiceEntityAlarmInstructionMenu from './service-entity-alarm-instruction-menu.vue';

export default {
  components: {
    ServiceEntityAlarmInstructionMenu,
  },
  props: {
    entity: {
      type: Object,
      default: () => ({}),
    },
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
