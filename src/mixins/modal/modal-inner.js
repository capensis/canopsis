import { createNamespacedHelpers } from 'vuex';

import modalMixin from './modal';

const { mapGetters: modalMapGetters } = createNamespacedHelpers('modal');

/**
 * @mixin
 */
export default {
  mixins: [modalMixin],
  computed: {
    ...modalMapGetters(['config']),
  },
};
