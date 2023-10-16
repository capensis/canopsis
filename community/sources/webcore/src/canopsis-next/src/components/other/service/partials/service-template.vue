<template>
  <c-compiled-template
    :template="modalTemplate"
    :context="templateContext"
  />
</template>

<script>
import Handlebars from 'handlebars';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

import ServiceEntitiesList from './service-entities-list.vue';

export default {
  components: {
    ServiceEntitiesList,
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
    widgetParameters: {
      type: Object,
      default: () => ({}),
    },
    pagination: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
  },
  computed: {
    templateContext() {
      return { entity: this.service };
    },

    modalTemplate() {
      return this.widgetParameters.modalTemplate || '';
    },
  },
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      return new Handlebars.SafeString(`
        <service-entities-list
          :service="service"
          :service-entities="serviceEntities"
          :widget-parameters="widgetParameters"
          :pagination="pagination"
          :total-items="totalItems"
          entity-name-field="${entityNameField}"
          @refresh="refreshEntities"
          @update:pagination="updatePagination"
        ></service-entities-list>
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
  methods: {
    updatePagination(pagination) {
      this.$emit('update:pagination', pagination);
    },

    refreshEntities() {
      this.$emit('refresh');
    },
  },
};
</script>
