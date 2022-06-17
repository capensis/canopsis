<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    v-tooltip(v-if="hasPbehaviorListAccess", left)
      template(#activator="{ on }")
        v-btn.pbehavior-modal-btn(
          v-on="on",
          small,
          dark,
          @click="showPbehaviorsListModal"
        )
          v-icon(small) list
        span {{ $t('modals.service.editPbehaviors') }}
    v-runtime-template(v-if="compiledTemplate && !pending", :template="compiledTemplate")
    div.float-clear
    c-table-pagination(
      v-if="!pending && totalItems > pagination.rowsPerPage && hasEntitiesHelper",
      :total-items="totalItems",
      :rows-per-page="pagination.rowsPerPage",
      :page="pagination.page",
      @update:page="updatePage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { PAGINATION_LIMIT } from '@/config';

import { CRUD_ACTIONS, MODALS, USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import ServiceEntitiesWrapper from './service-entities-wrapper.vue';

export default {
  components: {
    VRuntimeTemplate,
    ServiceEntitiesWrapper,
  },
  mixins: [authMixin],
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
  computed: {
    hasPbehaviorListAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.pbehaviorList);
    },

    hasEntitiesHelper() {
      return /{{(\s)?entities(.+)}}/.test(this.modalTemplate);
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
          entity-name-field="${entityNameField}"
          @add:action="addActionToQueue"
          @refresh="refreshEntities"
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

    refreshEntities() {
      this.$emit('refresh');
    },

    showPbehaviorsListModal() {
      this.$modals.show({
        name: MODALS.pbehaviorList,
        config: {
          pbehaviors: this.service.pbehaviors,
          entityId: this.service._id,
          onlyActive: true,
          availableActions: [CRUD_ACTIONS.create, CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
      });
    },

    updatePage(page) {
      this.$emit('update:pagination', { ...this.pagination, page });
    },

    updateRecordsPerPage(rowsPerPage) {
      this.$emit('update:pagination', { ...this.pagination, rowsPerPage, page: 1 });
    },
  },
};
</script>

<style lang="scss">
  .pbehavior-modal-btn {
    float: right;
  }
  .float-clear {
    clear: both;
  }
</style>
