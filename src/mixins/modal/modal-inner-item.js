import { createNamespacedHelpers } from 'vuex';

import ModalInnerMixin from './modal-inner';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

export default {
  mixins: [ModalInnerMixin],
  computed: {
    ...entitiesMapGetters(['getItem', 'getList']),

    item() {
      return this.getItem(this.config.itemType, this.config.itemId);
    },
  },
};
