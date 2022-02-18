<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.selectView.title') }}
    template(#text="")
      v-fade-transition
        v-layout(v-if="pending", justify-center)
          v-progress-circular(color="primary", indeterminate)
        v-layout(v-else)
          v-expansion-panel(dark)
            v-expansion-panel-content.secondary(v-for="group in groups", :key="group._id", ripple)
              template(#header="")
                div {{ group.title }}
              v-list.py-0.px-2.secondary
                v-list-tile.secondary.lighten-1(
                  v-for="view in group.views",
                  :key="view._id",
                  ripple,
                  @click="selectView(view._id)"
                )
                  v-list-tile-title.body-1 {{ view.title }}
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectView,
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
    async selectView(viewId) {
      if (this.config.action) {
        await this.config.action(viewId);
      }

      this.$modals.hide();
    },
  },
};
</script>
