<template lang="pug">
  v-menu(offset-y, @click.prevent.stop)
    v-btn(slot="activator", v-on="$listeners", depressed, small, light)
      v-icon {{ icon }}
    v-list
      service-entity-alarm-instruction-menu-item(
        v-for="assignedInstruction in assignedInstructions",
        :key="assignedInstruction._id",
        :assigned-instruction="assignedInstruction",
        @execute="$listeners.execute"
      )
</template>

<script>
import ServiceEntityAlarmInstructionMenuItem from './service-entity-alarm-instruction-menu-item.vue';

export default {
  components: {
    ServiceEntityAlarmInstructionMenuItem,
  },
  props: {
    assignedInstructions: {
      type: Array,
      default: () => [],
    },
    icon: {
      type: String,
      required: true,
    },
  },
};
</script>
