import { createNamespacedHelpers } from 'vuex';

import ModalMixin from './modal';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

export default {
  mixins: [ModalMixin],
  computed: {
    ...entitiesMapGetters(['getItem', 'getList']),

    item() {
      return this.getItem(this.config.itemType, this.config.itemId);
    },

    items() {
      return this.getList(this.config.itemType, this.config.itemsIds);
    },
  },
};
