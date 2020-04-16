<template lang="pug">
  div
    v-tooltip(v-if="hasPbehaviorListAccess", left)
      v-btn.pbehavior-modal-btn(
        slot="activator",
        small,
        dark,
        @click="showPbehaviorsListModal"
      )
        v-icon(small) edit
      span {{ $t('modals.watcher.editPbehaviors') }}
    v-runtime-template(:template="compiledTemplate")
    .float-clear
    v-layout.white(v-if="watchersMeta.total", align-center)
      v-flex(xs10)
        pagination(
          :page="watchersMeta.page",
          :limit="watchersMeta.limit",
          :total="watchersMeta.total",
          @input="updateQueryPage"
        )
      v-spacer
      v-flex(xs2)
        records-per-page(:value="watchersMeta.limit", @input="updateRecordsPerPage")
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { CRUD_ACTIONS, MODALS, USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';

import Pagination from '@/components/tables/pagination.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import WatcherEntity from './entity.vue';

export default {
  components: {
    VRuntimeTemplate,
    WatcherEntity,
    RecordsPerPage,
    Pagination,
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
    watchersMeta: {
      type: Object,
      required: true,
    },
    modalTemplate: {
      type: String,
      default: '',
    },
    entityTemplate: {
      type: String,
      default: '',
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
  },
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      return new Handlebars.SafeString(`
        <div class="mt-2" v-for="watcherEntity in watcherEntities" :key="watcherEntity._id">
          <watcher-entity
            :watcher-id="watcher.entity_id"
            :is-watcher-on-pbehavior="watcher.active_pb_watcher"
            :entity="watcherEntity"
            :template="entityTemplate"
            entity-name-field="${entityNameField}"
            @add:event="addEventToQueue"
          ></watcher-entity>
        </div>
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
          pbehaviors: this.watcher.watcher_pbehavior,
          entityId: this.watcher.entity_id,
          onlyActive: true,
          availableActions: [CRUD_ACTIONS.create, CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
      });
    },
    updateQueryPage(page) {
      this.$emit('change:page', page);
    },
    updateRecordsPerPage(limit) {
      this.$emit('change:limit', limit);
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
