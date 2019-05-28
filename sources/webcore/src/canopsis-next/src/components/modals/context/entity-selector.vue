<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
        v-btn(icon, dark, @click.native="hideModal")
          v-icon close
    v-card-text
      entities-single-selection-list(
      v-model="selectedItem",
      :filterPreparer="config.filterPreparer",
      :initialSearchingText="config.entityId"
      )
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import EntitiesSingleSelectionList from '@/components/other/context/entities-single-selection-list.vue';

export default {
  name: MODALS.contextEntitySelector,
  components: { EntitiesSingleSelectionList },
  mixins: [modalInnerMixin],
  data() {
    return {
      selectedItem: this.modal.config.entityId || null,
    };
  },
  computed: {
    title() {
      if (this.config.title) {
        return this.config.title;
      }

      return this.$t('modals.contextEntitySelector.title');
    },
  },
  methods: {
    submit() {
      if (this.config.action) {
        this.config.action(this.selectedItem);
      }

      this.hideModal();
    },
  },
};
</script>
