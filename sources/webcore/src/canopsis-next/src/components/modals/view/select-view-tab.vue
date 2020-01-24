<template lang="pug">
  modal-wrapper(data-test="selectViewTabModal")
    template(slot="title")
      span {{ $t('modals.selectViewTab.title') }}
    template(slot="text")
      v-expansion-panel(dark)
        v-expansion-panel-content.secondary(v-for="group in availableGroups", :key="group._id", ripple)
          template(slot="header")
            div {{ group.name }}
          v-expansion-panel.px-2(dark)
            v-expansion-panel-content.secondary.lighten-1(v-for="view in group.views", :key="view._id", ripple)
              template(slot="header")
                div {{ view.title }}
              v-list.pa-0
                v-list-tile.secondary.lighten-2(
                  v-for="tab in view.tabs",
                  :key="tab._id",
                  ripple,
                  @click="selectTab(tab._id, view._id)"
                )
                  v-list-tile-title.body-1.pl-4 {{ tab.title }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewsGroupsMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectViewTab,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesViewsGroupsMixin,
    rightsEntitiesGroupMixin,
  ],
  methods: {
    async selectTab(tabId, viewId) {
      if (this.config.action) {
        await this.config.action({ tabId, viewId });
      }

      this.$modals.hide();
    },
  },
};
</script>
