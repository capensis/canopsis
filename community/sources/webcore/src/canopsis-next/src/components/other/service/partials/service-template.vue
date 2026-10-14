<template>
  <service-entities-list
    v-if="isEntitiesOnlyTemplate"
    :service="service"
    :service-entities="serviceEntities"
    :widget-parameters="widgetParameters"
    :options="options"
    :total-items="totalItems"
    :actions-requests="actionsRequests"
    :entity-name-field="entitiesNameField"
    @refresh="refreshEntities"
    @update:options="updateOptions"
    @add:action="addAction"
  />
  <c-compiled-template
    v-else
    :template="modalTemplate"
    :context="templateContext"
    :template-props="templateProps"
  />
</template>

<script>
import Handlebars from 'handlebars';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

import { handlebarsTagsHelperMixin } from '@/mixins/handlebars/tags-helper';

import ServiceEntitiesList from './service-entities-list.vue';

const ENTITIES_ONLY_TEMPLATE_REGEX = /^\{\{\s*entities(?:\s+name="[^"]+")?\s*\}\}$/;

export default {
  components: {
    ServiceEntitiesList,
  },

  mixins: [handlebarsTagsHelperMixin],

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
  emits: ['add:action', 'refresh', 'update:options'],
  computed: {
    templateContext() {
      return { entity: this.service };
    },

    templateProps() {
      return {
        entity: this.service,
        service: this.service,
        serviceEntities: this.serviceEntities,
        widgetParameters: this.widgetParameters,
        options: this.options,
        totalItems: this.totalItems,
        actionsRequests: this.actionsRequests,
      };
    },

    modalTemplate() {
      return this.widgetParameters.modalTemplate || '';
    },

    normalizedModalTemplate() {
      return this.modalTemplate.trim().replace(/\s+/g, ' ');
    },

    isEntitiesOnlyTemplate() {
      const { normalizedModalTemplate } = this;

      if (!normalizedModalTemplate) {
        return true;
      }

      return ENTITIES_ONLY_TEMPLATE_REGEX.test(normalizedModalTemplate);
    },

    entitiesNameField() {
      const match = this.modalTemplate.match(/\{\{\s*entities(?:\s+name="([^"]+)")?\s*\}\}/);

      return match?.[1] || 'entity._id';
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
  beforeUnmount() {
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
