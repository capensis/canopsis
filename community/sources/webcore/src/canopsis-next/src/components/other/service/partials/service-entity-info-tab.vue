<template>
  <div>
    <v-layout
      v-if="actions.length"
      class="pa-2"
      align-center
      wrap
    >
      <div>{{ $t('common.actionsLabel') }}:</div>
      <service-entity-actions
        :actions="actions"
        :entity="entity"
        :assigned-instructions="entity.assigned_instructions"
        @apply="$listeners.apply"
        @execute="$listeners.execute"
      />
    </v-layout>
    <service-entity-template
      :entity="entity"
      :template="template"
    />
  </div>
</template>

<script>
import ServiceEntityTemplate from './service-entity-template.vue';
import ServiceEntityActions from './service-entity-actions.vue';

export default {
  components: {
    ServiceEntityActions,
    ServiceEntityTemplate,
  },
  props: {
    entity: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
      default: '',
    },
    actions: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
