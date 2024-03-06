<template lang="pug">
  c-runtime-template(:template="compiledTemplate")
</template>

<script>
import Handlebars from 'handlebars';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

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
    actionsRequests: {
      type: Array,
      default: () => [],
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.modalTemplate, { entity: this.service });

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
  computed: {
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
          :actions-requests="actionsRequests"
          entity-name-field="${entityNameField}"
          @refresh="refreshEntities"
          @apply:action="applyAction"
          @update:pagination="updatePagination"
        ></service-entities-list>
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
  methods: {
    applyAction(action) {
      this.$emit('apply:action', action);
    },

    updatePagination(pagination) {
      this.$emit('update:pagination', pagination);
    },

    refreshEntities() {
      this.$emit('refresh');
    },
  },
};
</script>
