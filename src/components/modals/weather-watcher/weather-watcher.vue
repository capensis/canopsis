<template lang="pug">
  v-card
    v-card-title
      span.headline {{ watcher.display_name }}
    v-divider
    v-card-text
      div {{ $t('modals.watcherData.name') }}: {{ watcher.display_name }}
      div {{ $t('modals.watcherData.org') }}: {{ watcher.org }}
      div.entities-list(v-if="!watchedEntitiesPending", v-for="watchedEntity in watchedEntities")
        entity-data(:watcher="watcher", :watchedEntity="watchedEntity")
</template>

<script>
import EntityData from '@/components/modals/weather-watcher/entity-data.vue';
import watcherMixin from '@/mixins/watcher';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

export default {
  name: MODALS.weatherWatcher,
  components: {
    EntityData,
  },
  mixins: [
    watcherMixin,
    modalInnerMixin,
  ],
  computed: {
    watcher() {
      return this.getWatcher(this.config.watcherId);
    },
  },
  mounted() {
    this.fetchWatchedEntities({ params: {}, watcherId: this.config.watcherId });
  },
};
</script>

<style scoped>
  .entities-list {
    margin-top: 20px;
  }
</style>
