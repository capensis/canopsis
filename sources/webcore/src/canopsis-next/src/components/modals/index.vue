<template lang="pug">
  div
    modal-base(:modals="modals")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import ModalBase from './modal-base.vue';

const { mapGetters: modalMapGetters } = createNamespacedHelpers('modal');

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
};
</script>
