<template>
  <div>
    <c-layout
      v-if="actions.length"
      align="center"
      class="pa-2 gap-2"
      wrap
    >
      <div>{{ $t('common.actionsLabel') }}:</div>
      <service-entity-actions
        :actions="actions"
        :entity="entity"
        :assigned-instructions="entity.assigned_instructions"
        @apply="$emit('apply', $event)"
        @execute="$emit('execute', $event)"
      />
    </c-layout>
    <service-entity-template
      :entity="entity"
      :template="template"
    />
  </div>
</template>

<script>
import ServiceEntityActions from './service-entity-actions.vue';
import ServiceEntityTemplate from './service-entity-template.vue';

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
  emits: ['apply', 'execute'],
};
</script>
