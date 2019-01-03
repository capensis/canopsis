import { createNamespacedHelpers } from 'vuex';
import first from 'lodash/first';

import modalInnerMixin from './inner';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

/**
 * @mixin
 */
export default {
  mixins: [modalInnerMixin],
  computed: {
    ...entitiesMapGetters(['getList']),

    /**
     * Function returns first item from items array
     *
     * @returns {Object|{}}
     */
    firstItem() {
      return first(this.items) || {};
    },

    /**
     * Function returns items by entity type and entity ids
     *
     * @returns {Array}
     */
    items() {
      return this.getList(this.config.itemsType, this.config.itemsIds);
    },
  },
};
