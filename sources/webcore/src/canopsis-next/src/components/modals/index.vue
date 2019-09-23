<template lang="pug">
  div
    modal-base(:modals="modals")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import ModalBase from './modal-base.vue';

const {
  mapGetters: modalMapGetters,
} = createNamespacedHelpers('modal');

/**
 * Wrapper for all modal windows
 */
export default {
  components: {
    ModalBase,
  },
  computed: {
    ...modalMapGetters(['modals']),
  },
  watch: {
    $route: {
      handler() {
        if (this.modals && this.modals.length) {
          this.modals.map(modal => this.$modals.hide({ id: modal.id }));
        }
      },
      deep: true,
    },
  },
};
</script>
