<template lang="pug">
  div
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

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
      return `<div>${compile(this.modalTemplate, { entity: this.watcher })}</div>`;
    },
  },
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      return new Handlebars.SafeString(`
        <div class="mt-2" v-for="watcherEntity in watcherEntities">
          <watcher-entity
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
  },
};
</script>
