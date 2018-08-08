<template lang="pug">
  v-card
    v-card-title
      span.headline {{ watcher.display_name }}
    v-divider
    v-card-text
      div {{ $t('modals.weatherWatcher.name') }}: {{ watcher.display_name }}
      div {{ $t('modals.weatherWatcher.org') }}: {{ watcher.org }}
      div.entities-list(v-if="!watchedEntitiesPending", v-for="watchedEntity in watchedEntities")
        entity(:entity="watchedEntity")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import Entity from '@/components/modals/weather-watcher/entity.vue';
import watcherMixin from '@/mixins/watcher';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

const { mapGetters } = createNamespacedHelpers('weatherWatcher');


export default {
  name: MODALS.weatherWatcher,
  components: {
    Entity,
  },
  mixins: [
    watcherMixin,
    modalInnerMixin,
  ],
  computed: {
    ...mapGetters({
      getWeatherWatcher: 'getItem',
    }),
    watcher() {
      return this.getWeatherWatcher(this.config.watcherId);
    },
  },
  mounted() {
    this.fetchWatchedEntities({ watcherId: this.config.watcherId });
  },
};
</script>

<style scoped>
  .entities-list {
    margin-top: 20px;
  }
</style>
