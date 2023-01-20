<template lang="pug">
  div
    v-layout.pa-2(v-if="actions.length", align-center, wrap)
      div {{ $t('common.actionsLabel') }}:
      service-entity-actions(
        :actions="actions",
        :entity="entity",
        :assigned-instructions="entity.assigned_instructions",
        @apply="$listeners.apply",
        @execute="$listeners.execute"
      )
    service-entity-template(:entity="entity", :template="template")
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
