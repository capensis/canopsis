<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.watcherLists.title') }}
        v-btn(icon, dark, @click.native="hideModal")
          v-icon close
    v-card-text
      entities-list-widget(:widget="config.widget")
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import queryMixin from '@/mixins/query';

import EntitiesListWidget from '@/components/other/context/entities-list.vue';

export default {
  name: MODALS.watchersList,
  components: { EntitiesListWidget },
  mixins: [modalInnerMixin, queryMixin],
  created() {
    this.mergeQuery({
      id: this.config.widget._id,
      query: this.config.query,
    });
  },
};
</script>
