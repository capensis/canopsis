import { createNamespacedHelpers } from 'vuex';

import ModalInnerMixin from './modal-inner';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

/**
 * @mixin
 */
export default {
  mixins: [ModalInnerMixin],
  computed: {
    ...entitiesMapGetters(['getItem', 'getList']),

    /**
     * Function returns item by entity type and entity id
     *
     * @returns {Object}
     */
    item() {
      return this.getItem(this.config.itemType, this.config.itemId);
    },
  },
};
