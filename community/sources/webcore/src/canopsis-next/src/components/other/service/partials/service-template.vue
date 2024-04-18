<template>
  <c-compiled-template
    :template="modalTemplate"
    :context="templateContext"
    :template-props="templateProps"
  />
</template>

<script>
import Handlebars from 'handlebars';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

import { entityHandlebarsTagsHelper } from '@/mixins/widget/handlebars/entity-tags-helper';

import ServiceEntitiesList from './service-entities-list.vue';

export default {
  components: {
    // eslint-disable-next-line vue/no-unused-components
    ServiceEntitiesList,
  },
  mixins: [entityHandlebarsTagsHelper],
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
    options: {
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
  computed: {
    templateContext() {
      return { entity: this.service };
    },

    templateProps() {
      return { entity: this.service };
    },

    modalTemplate() {
      return this.widgetParameters.modalTemplate || '';
    },
  },
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      /**
       * For new properties and events you must put it for sanitizer
       */
      return new Handlebars.SafeString(`
        <service-entities-list
          :service="service"
          :service-entities="serviceEntities"
          :widget-parameters="widgetParameters"
          :options="options"
          :total-items="totalItems"
          :actions-requests="actionsRequests"
          entity-name-field="${entityNameField}"
          @refresh="refreshEntities"
          @update:options="updateOptions"
          @add:action="addAction"
        ></service-entities-list>
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
  methods: {
    addAction(action) {
      this.$emit('add:action', action);
    },

    updateOptions(pagination) {
      this.$emit('update:options', pagination);
    },

    refreshEntities(immediate = false) {
      this.$emit('refresh', immediate);
    },
  },
};
</script>
