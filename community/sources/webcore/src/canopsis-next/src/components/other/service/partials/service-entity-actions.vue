<template>
  <c-layout class="gap-4 ml-2" wrap>
    <c-tooltip
      v-for="action in actions"
      :key="action.type"
      location="top"
    >
      <template #activator="{ props: tooltipProps }">
        <span v-bind="tooltipProps">
          <service-entity-alarm-instruction-menu
            v-if="action.type === $constants.WEATHER_ACTIONS_TYPES.executeInstruction"
            :icon="action.icon"
            :entity="entity"
            :assigned-instructions="assignedInstructions"
            @execute="$emit('execute', $event)"
          />
          <v-btn
            v-else
            :disabled="action.disabled"
            :loading="action.loading"
            class="ml-2"
            size="small"
            variant="flat"
            @click.stop="$emit('apply', action)"
          >
            <v-icon :icon="action.icon" />
          </v-btn>
        </span>
      </template>
      <span>{{ $t(`serviceWeather.actions.${action.type}`) }}</span>
    </c-tooltip>
  </c-layout>
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
  emits: ['apply', 'execute'],
};
</script>
