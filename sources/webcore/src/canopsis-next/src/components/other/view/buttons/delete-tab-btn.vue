<template lang="pug">
  v-btn(
    data-test="deleteTab",
    small,
    flat,
    icon,
    @click.prevent="showDeleteTabModal(tab)"
  )
    v-icon(small) delete
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [
    modalMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
    view: {
      type: Object,
      required: true,
    },
    updateViewMethod: {
      type: Function,
      default: () => {},
    },
  },
  methods: {
    showDeleteTabModal(tab) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.deleteTab(tab._id),
        },
      });
    },

    deleteTab(tabId) {
      const view = {
        ...this.view,
        tabs: this.view.tabs.filter(viewTab => viewTab._id !== tabId),
      };

      return this.updateViewMethod(view);
    },
  },
};
</script>
