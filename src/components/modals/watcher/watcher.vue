<template lang="pug">
  v-card
    v-card-title
      span.headline {{ watcher.display_name }}
    v-divider
    v-card-text
      v-layout(v-for="attribute in Object.keys(attributes)", row, wrap)
        v-flex.text-md-right(xs3)
          b {{ $t(`modals.weatherWatcher.${attribute}`) }}:
        v-flex.pl-2(xs9)
          span {{ attributes[attribute] }}
      div(v-if="!watcherEntitiesPending")
        div.mt-4(v-for="watcherEntity in watcherEntities")
          watcher-entity(:entity="watcherEntity")
</template>

<script>
import pick from 'lodash/pick';
import mapValues from 'lodash/mapValues';

import { MODALS } from '@/constants';
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
