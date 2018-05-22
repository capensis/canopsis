import { createNamespacedHelpers } from 'vuex';

import ModalMixin from './modal';

const { mapGetters: modalMapGetters } = createNamespacedHelpers('modal');

export default {
  mixins: [ModalMixin],
  computed: {
    ...modalMapGetters(['config']),
  },
};
