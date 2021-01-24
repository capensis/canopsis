<template lang="pug">
  div
    div.mt-2(v-for="watcherEntity in watcherEntities", :key="watcherEntity.key")
      watcher-entity(
        :watcher-id="watcher._id",
        :is-watcher-on-pbehavior="watcher.active_pb_watcher",
        :entity="watcherEntity",
        :template="entityTemplate",
        :entity-name-field="entityNameField",
        @add:event="$emit('add:event', $event)"
      )
</template>

<script>
import WatcherEntity from './entity.vue';

export default {
  components: {
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
    entityNameField: {
      type: String,
      default: 'name',
    },
    entityTemplate: {
      type: String,
      default: '',
    },
  },
};
</script>
