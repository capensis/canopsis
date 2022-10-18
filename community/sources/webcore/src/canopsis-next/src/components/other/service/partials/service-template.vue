<template lang="pug">
  div.position-relative
    v-layout.pa-4(v-if="pending", justify-center)
      v-progress-circular(color="primary", indeterminate)
    v-runtime-template(v-else-if="compiledTemplate", :template="compiledTemplate")
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { PAGINATION_LIMIT } from '@/config';

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
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      return new Handlebars.SafeString(`
        <service-entities-wrapper
          :service="service"
          :service-entities="serviceEntities"
          :widget-parameters="widgetParameters"
          :pagination="pagination"
          :total-items="totalItems"
          entity-name-field="${entityNameField}"
          @add:action="addActionToQueue"
          @refresh="refreshEntities"
          @update:pagination="updatePagination"
        ></service-entities-wrapper>
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
  methods: {
    addActionToQueue(event) {
      this.$emit('add:action', event);
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
