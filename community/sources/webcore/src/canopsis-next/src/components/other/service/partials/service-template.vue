<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    v-runtime-template(:key="templateKey", :template="compiledTemplate")
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { PAGINATION_LIMIT } from '@/config';

import uid from '@/helpers/uid';
import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import PbehaviorsSimpleList from '@/components/other/pbehavior/partials/pbehaviors-simple-list.vue';

import ServiceEntitiesWrapper from './service-entities-wrapper.vue';

export default {
  components: {
    PbehaviorsSimpleList,
    VRuntimeTemplate,
    ServiceEntitiesWrapper,
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
    itemsPerPage: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
    pagination: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
      required: false,
    },
    unavailableEntitiesAction: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      templateKey: uid(),
    };
  },
  asyncComputed: {
    modalTemplate() {
      return this.widgetParameters.modalTemplate || '';
    },

    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.modalTemplate, { entity: this.service });

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
  watch: {
    service: 'generateTemplateKey',
    serviceEntities: 'generateTemplateKey',
  },
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      return new Handlebars.SafeString(`
        <service-entities-wrapper
          :service="service"
          :service-entities="serviceEntities"
          :unavailable-entities-action="unavailableEntitiesAction"
          :widget-parameters="widgetParameters"
          :pagination="pagination"
          :total-items="totalItems"
          entity-name-field="${entityNameField}"
          @refresh="refreshEntities"
          @apply:action="applyAction"
          @remove:unavailable="removeUnavailable"
          @update:pagination="updatePagination"
        ></service-entities-wrapper>
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
  methods: {
    generateTemplateKey() {
      this.templateKey = uid();
    },

    applyAction(event) {
      this.$emit('apply:action', event);
    },

    updatePagination(pagination) {
      this.$emit('update:pagination', pagination);
    },

    removeUnavailable(entity) {
      this.$emit('remove:unavailable', entity);
    },

    refreshEntities() {
      this.$emit('refresh');
    },
  },
};
</script>
