<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.watcherLists.title') }}
        v-btn(icon, dark, @click.native="hideModal")
          v-icon close
    v-card-text
      entities-list-widget(:widget="config.widget")
        template(slot="item-selector", slot-scope="{ item }")
          v-radio-group(v-model="selected", hide-details)
            v-radio(:value="item._id")
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import queryMixin from '@/mixins/query';

import EntitiesListWidget from '@/components/other/context/entities-list.vue';

export default {
  name: MODALS.contextEntitiesList,
  components: { EntitiesListWidget },
  mixins: [modalInnerMixin, queryMixin],
  data() {
    return {
      selectedItem: null,
    };
  },
  created() {
    this.mergeQuery({
      id: this.config.widget._id,
      query: this.config.query,
    });
  },
};
</script>
