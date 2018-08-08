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
      div.mt-4(v-if="!watchedEntitiesPending", v-for="watchedEntity in watchedEntities")
        weather-watcher-entity(:entity="watchedEntity")
</template>

<script>
import pick from 'lodash/pick';
import mapValues from 'lodash/mapValues';

import { MODALS } from '@/constants';
import WeatherWatcherEntity from '@/components/modals/weather-watcher/entity.vue';

import entitiesWatcherMixin from '@/mixins/entities/watcher';
import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.weatherWatcher,
  components: {
    WeatherWatcherEntity,
  },
  mixins: [
    entitiesWatcherMixin,
    modalInnerMixin,
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

    //    this.fetchWatchedEntities({ watcherId: this.config.watcherId });
  },
};
</script>
