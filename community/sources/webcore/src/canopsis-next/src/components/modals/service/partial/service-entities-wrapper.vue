<template lang="pug">
  div
    div.mt-2(v-for="serviceEntity in serviceEntities", :key="serviceEntity.key")
      service-entity(
        :service-id="service._id",
        :entity="serviceEntity",
        :entity-name-field="entityNameField",
        :widget-parameters="widgetParameters",
        @add:event="$listeners['add:event']",
        @refresh="$listeners.refresh"
      )
</template>

<script>
import ServiceEntity from './service-entity.vue';

export default {
  components: {
    ServiceEntity,
  },
  props: {
    service: {
      type: Object,
      required: true,
    },
    serviceEntities: {
      type: Array,
      default: () => [],
    },
    entityNameField: {
      type: String,
      default: 'name',
    },
    widgetParameters: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
