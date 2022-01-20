<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.selectView.title') }}
    template(#text="")
      v-fade-transition
        c-progress-overlay(v-if="pending")
        v-expansion-panel(v-else, dark)
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
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('view/group');

export default {
  name: MODALS.selectView,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
  ],
  data() {
    return {
      pending: true,
      groups: [],
    };
  },
  async mounted() {
    this.pending = true;

    const { data } = await this.fetchGroupsListWithoutStore({
      params: {
        limit: MAX_LIMIT,
        page: 1,
        with_views: true,
        with_flags: true,
      },
    });

    this.groups = data;
    this.pending = false;
  },
  methods: {
    ...mapActions({
      fetchGroupsListWithoutStore: 'fetchListWithoutStore',
    }),

    async selectView(viewId) {
      if (this.config.action) {
        await this.config.action(viewId);
      }

      this.$modals.hide();
    },
  },
};
</script>
