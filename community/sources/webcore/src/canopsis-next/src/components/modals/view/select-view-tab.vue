<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.selectViewTab.title') }}
    template(#text="")
      v-fade-transition
        v-layout(v-if="pending", justify-center)
          v-progress-circular(color="primary", indeterminate)
        v-layout(v-else)
          v-expansion-panel(dark)
            v-expansion-panel-content.secondary(
              v-for="group in groups",
              :key="group._id",
              ripple
            )
              template(#header="")
                div {{ group.title }}
              v-expansion-panel.px-2(dark)
                v-expansion-panel-content.secondary.lighten-1(
                  v-for="view in group.views",
                  :key="view._id",
                  ripple
                )
                  template(#header="")
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectViewTab,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesViewGroupMixin,
  ],
  data() {
    return {
      pending: true,
    };
  },
  async mounted() {
    this.pending = true;

    await this.fetchAllGroupsListWithWidgetsWithCurrentUser();

    this.pending = false;
  },
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
