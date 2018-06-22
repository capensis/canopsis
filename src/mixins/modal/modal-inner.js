// LIBS
import { createNamespacedHelpers } from 'vuex';
// MIXINS
import modalMixin from '@/mixins/modal/modal';

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
