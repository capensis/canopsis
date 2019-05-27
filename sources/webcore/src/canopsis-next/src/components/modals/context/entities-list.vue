<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.watcherLists.title') }}
        v-btn(icon, dark, @click.native="hideModal")
          v-icon close
    v-card-text
      context-general-list(
      v-model="selectedItem",
      :filterPreparer="filterPreparer",
      initialSearchingText="watcher_9b55e9cb-e050-4c20-ac74-1df91c52e038",
      single
      )
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, CONTEXT_ENTITIES_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import queryMixin from '@/mixins/query';

import { getContextSearchByText } from '@/helpers/widget-search';

import ContextGeneralList from '@/components/other/context/context-general-list.vue';

export default {
  name: MODALS.contextEntitiesList,
  components: { ContextGeneralList },
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
  methods: {
    filterPreparer(text) {
      return {
        $and: [
          getContextSearchByText(text, ['_id']),

          {
            type: CONTEXT_ENTITIES_TYPES.watcher,
          },
        ],
      };
    },
  },
};
</script>
