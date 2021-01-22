<template lang="pug">
  div
    v-tooltip(v-if="hasPbehaviorListAccess", left)
      v-btn.pbehavior-modal-btn(
        slot="activator",
        small,
        dark,
        @click="showPbehaviorsListModal"
      )
        v-icon(small) list
      span {{ $t('modals.watcher.editPbehaviors') }}
    v-runtime-template(v-if="compiledTemplate", :template="compiledTemplate")
    .float-clear
    c-table-pagination(
      v-if="totalItems > pagination.rowsPerPage && hasEntitiesHelper",
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

import { CRUD_ACTIONS, MODALS, USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import WatcherEntitiesWrapper from './entities-wrapper.vue';

export default {
  components: {
    VRuntimeTemplate,
    WatcherEntitiesWrapper,
  },
  mixins: [authMixin],
  props: {
    watcher: {
      type: Object,
      required: true,
    },
    watcherEntities: {
      type: Array,
      default: () => [],
    },
    modalTemplate: {
      type: String,
      default: '',
    },
    entityTemplate: {
      type: String,
      default: '',
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
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.modalTemplate, { entity: this.watcher });

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
  computed: {
    hasPbehaviorListAccess() {
      return this.checkAccess(USERS_RIGHTS.business.weather.actions.pbehaviorList);
    },

    hasEntitiesHelper() {
      return /{{(\s)?entities(.+)}}/.test(this.modalTemplate);
    },
  },
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      return new Handlebars.SafeString(`
        <watcher-entities-wrapper
          :watcher="watcher"
          :watcher-entities="watcherEntities"
          :entity-template="entityTemplate"
          entity-name-field="${entityNameField}"
          @add:event="addEventToQueue"
        />
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
  methods: {
    addEventToQueue(event) {
      this.$emit('add:event', event);
    },

    showPbehaviorsListModal() {
      this.$modals.show({
        name: MODALS.pbehaviorList,
        config: {
          pbehaviors: this.watcher.pbehaviors,
          entityId: this.watcher._id,
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
