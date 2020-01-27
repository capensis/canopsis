<template lang="pug">
  div
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { CRUD_ACTIONS, MODALS } from '@/constants';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import WatcherEntity from './entity.vue';

export default {
  components: {
    VRuntimeTemplate,
    WatcherEntity,
  },
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
  },
  computed: {
    compiledTemplate() {
      return `
        <div>
          <v-tooltip left>
            <v-btn small dark slot="activator" class="pbehavior-modal-btn" @click="showPbehaviorsListModal">
              <v-icon small>edit</v-icon>
            </v-btn>
            <span>${this.$t('modals.watcher.editPbehaviors')}</span>
          </v-tooltip>
          ${compile(this.modalTemplate, { entity: this.watcher })}
          <div class="float-clear"></div>
        </div>
      `;
    },
  },
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      return new Handlebars.SafeString(`
        <div class="mt-2" v-for="watcherEntity in watcherEntities" :key="watcherEntity._id">
          <watcher-entity
          :watcherId="watcher.entity_id"
          :isWatcherOnPbehavior="watcher.active_pb_watcher"
          :entity="watcherEntity"
          :template="entityTemplate"
          entityNameField="${entityNameField}"
          @addEvent="addEventToQueue"></watcher-entity>
        </div>
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
  methods: {
    addEventToQueue(event) {
      this.$emit('addEvent', event);
    },
    showPbehaviorsListModal() {
      this.$modals.show({
        name: MODALS.pbehaviorList,
        config: {
          pbehaviors: this.watcher.watcher_pbehavior,
          entityId: this.watcher.entity_id,
          onlyActive: true,
          showAddButton: true,
          availableActions: [CRUD_ACTIONS.delete, CRUD_ACTIONS.update],
        },
      });
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
