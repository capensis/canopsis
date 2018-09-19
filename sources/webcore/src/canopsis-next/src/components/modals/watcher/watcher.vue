<template lang="pug">
  v-card
    v-card-title
      span.headline {{ watcher.display_name }}
    v-divider
    v-card-text
      div(v-html="compiledTemplate")
      div(v-if="!watcherEntitiesPending")
        div.mt-2(v-for="watcherEntity in watcherEntities")
          watcher-entity(:entity="watcherEntity", :template="config.entityTemplate")
</template>

<script>
import pick from 'lodash/pick';
import mapValues from 'lodash/mapValues';

import { MODALS } from '@/constants';
import compile from '@/helpers/handlebars';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';
import modalInnerMixin from '@/mixins/modal/modal-inner';

import WatcherEntity from './partial/entity.vue';


export default {
  name: MODALS.watcher,
  components: {
    WatcherEntity,
  },
  mixins: [
    modalInnerMixin,
    entitiesWatcherMixin,
    entitiesWatcherEntityMixin,
  ],
  data() {
    return {
      attributes: {},
    };
  },
  computed: {
    watcher() {
      return this.getWatcher(this.config.watcherId);
    },
    compiledTemplate() {
      return compile(this.config.modalTemplate, { watcher: this.watcher });
    },
  },
  mounted() {
    const infoAttributes = mapValues(pick(this.watcher.infos, [
      'application_crit_label',
      'product_line',
      'service_period',
      'isInCarat',
      'application_label',
      'target_platform',
    ]), v => v.value);

    this.attributes = {
      org: this.watcher.org,
      ...infoAttributes,
    };

    this.fetchWatcherEntitiesList({ watcherId: this.config.watcherId });
  },
};
</script>
